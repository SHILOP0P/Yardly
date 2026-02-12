package pgrepo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
    "github.com/jackc/pgx/v5/pgconn"


	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/SHILOP0P/Yardly/backend/internal/auth"
	"github.com/SHILOP0P/Yardly/backend/internal/user"
)

type Repo struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool)*Repo{
	return &Repo{pool: pool}
}
func (r *Repo) CreateWithProfile(ctx context.Context, u *user.User, p *user.Profile) error{
	tx, err := r.pool.Begin(ctx)
	if err != nil{
		log.Println("error begin:", err)
		return err
	}
	defer tx.Rollback(ctx)
	const uq = `
	INSERT INTO users (email, password_hash, role)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, updated_at
	`
	err = tx.QueryRow(ctx, uq, u.Email, u.PasswordHash,
		u.Role,
	).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
    	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
        return user.ErrEmailTaken
    }
    return fmt.Errorf("insert user: %w", err)
	}

	const pq = `
	INSERT INTO user_profiles (user_id, first_name, last_name)
	VALUES ($1, $2, $3)
	`
	_, err = tx.Exec(ctx, pq,
		u.ID,
		p.FirstName,
		p.LastName,
	)
	if err != nil{
		log.Println("insert profile error:", err)
		return fmt.Errorf("insert profile: %w", err)
	}
	return tx.Commit(ctx)
}

func (r *Repo) GetByEmail(ctx context.Context, email string) (user.User, error){
	const q = `
	SELECT id, email, password_hash, role, banned_at, ban_expires_at, ban_reason, created_at, updated_at
	FROM users
	WHERE email = $1
	`
	var u user.User
	err := r.pool.QueryRow(ctx, q, email).Scan(
		&u.ID,
		&u.Email,
		&u.PasswordHash,
		&u.Role,
		&u.BannedAt, 
		&u.BanExpiresAt,
		&u.BanReason,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user.User{}, user.ErrNotFound
		}
		return user.User{}, err
	}
	return u, nil
}

func (r *Repo) GetByID(ctx context.Context, id int64)(user.User, user.Profile, error){
	const q = `
	SELECT
		u.id, u.email, u.role, u.banned_at, u.ban_expires_at, u.ban_reason, u.created_at, u.updated_at,
		p.first_name, p.last_name, p.birth_date, p.gender, p.avatar_url, p.updated_at
	FROM users u
	JOIN user_profiles p ON p.user_id = u.id
	WHERE u.id = $1
	`
	var u user.User
	var p user.Profile
	err:=r.pool.QueryRow(ctx, q, id).Scan(
		&u.ID, &u.Email, &u.Role, &u.BannedAt, &u.BanExpiresAt, &u.BanReason, &u.CreatedAt, &u.UpdatedAt,
		&p.FirstName, &p.LastName, &p.BirthDate, &p.Gender, &p.AvatarURL, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user.User{}, user.Profile{}, user.ErrNotFound
		}
		return user.User{}, user.Profile{}, err
	}
	p.UserID = u.ID
	return u, p, nil
}

func (r *Repo) Authenticate(ctx context.Context, email, password string) (auth.AuthUser, error) {
	u, err := r.GetByEmail(ctx, email)
	if err != nil {
		return auth.AuthUser{}, err
	}

	// тут подставь свою реальную проверку пароля:
	// например user.CheckPassword(u.PasswordHash, password)
	if ok := user.CheckPassword(u.PasswordHash, password); !ok {
		// важно: не палим что email существует
		return auth.AuthUser{}, user.ErrNotFound
	}

	now := time.Now().UTC()
	banned := isBanned(u.BannedAt, u.BanExpiresAt, now)

	return auth.AuthUser{
		ID:     u.ID,
		Role:   auth.Role(u.Role),
		Banned: banned,
	}, nil
}


func (r *Repo) GetByIDForAuth(ctx context.Context, id int64) (auth.AuthUser, error) {
	const q = `
	SELECT id, role, banned_at, ban_expires_at
	FROM users
	WHERE id = $1
	`

	var userID int64
	var role user.Role
	var bannedAt *time.Time
	var bannedExpiresAt *time.Time
	var banned bool

	err := r.pool.QueryRow(ctx, q, id).Scan(&userID, &role, &bannedAt, &bannedExpiresAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return auth.AuthUser{}, user.ErrNotFound
		}
		return auth.AuthUser{}, err
	}

	now := time.Now().UTC()

	banned = isBanned(bannedAt, bannedExpiresAt, now)

	return auth.AuthUser{
		ID:     userID,
		Role:   auth.Role(role),
		Banned: banned,
	}, nil
}


func isBanned(bannedAt, banExpiresAt *time.Time, now time.Time) bool {
	// Бан никогда не выдавался
	if bannedAt == nil {
		return false
	}

	// Бан выдан и не имеет срока окончания → перманентный
	if banExpiresAt == nil {
		return true
	}

	// Временный бан: активен, пока now < banExpiresAt
	return now.Before(*banExpiresAt)
}
