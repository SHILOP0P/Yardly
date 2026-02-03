package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"time"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshRepo struct{
	pool *pgxpool.Pool
	secret []byte
}

func NewRefreshRepo(pool *pgxpool.Pool, secret string) *RefreshRepo{
	return &RefreshRepo{pool: pool, secret: []byte(secret)}
}

func (r *RefreshRepo) hashToken(raw string)[]byte{
	m := hmac.New(sha256.New, r.secret)
	_, _ = m.Write([]byte(raw))
	return m.Sum(nil)
}

func (r *RefreshRepo) Create(ctx context.Context, userID int64, rawToken string, expiresAt time.Time) error {
	const q = `
	INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
	VALUES ($1, $2, $3)
	`
	_, err := r.pool.Exec(ctx, q, userID, r.hashToken(rawToken), expiresAt)
	return err
}

func (r *RefreshRepo) Consume(ctx context.Context, rawToken string, now time.Time)(int64, error){
	const q = `
	UPDATE refresh_tokens
	SET revoked_at = $2
	WHERE token_hash = $1
	AND revoked_at IS NULL
	AND expires_at > $2
	RETURNING user_id
	`
	var userID int64
	err :=r.pool.QueryRow(ctx, q, r.hashToken(rawToken), now).Scan(&userID)
	return userID, err
}

func (r *RefreshRepo) Rotate(ctx context.Context, oldRaw, newRaw string, newExpiresAt, now time.Time, grace time.Duration)(int64, error){
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err!=nil{
		return 0, err
	}
	defer func(){
		_=tx.Rollback(ctx)
	}()

	const selQ = `
	SELECT id, user_id, revoked_at, expires_at, replaced_by
    FROM refresh_tokens
    WHERE token_hash = $1
    LIMIT 1
    FOR UPDATE
	`
	var (
		oldID int64
		userID int64
		revokedAt *time.Time
		expiresAt time.Time
		replaceBy *int64
	)

	err = tx.QueryRow(ctx, selQ, r.hashToken(oldRaw)).Scan(&oldID, &userID, &revokedAt, &expiresAt, &replaceBy)
	if err == pgx.ErrNoRows{
		return 0, ErrInvalidRefresh
	}
	if err!=nil{
		log.Println("Fail selecting for update")
		return 0, err
	}

	if revokedAt != nil{
		if replaceBy!=nil && now.Sub(*revokedAt)<=grace{
			return 0, ErrRefreshAlreadyRotated
		}
		        const revokeAllQ = `
        UPDATE refresh_tokens
        SET revoked_at = $2
        WHERE user_id = $1
          AND revoked_at IS NULL
        `
        if _, err := tx.Exec(ctx, revokeAllQ, userID, now); err != nil {
			log.Println("Fail in revoking all")
            return 0, err
        }
        if err := tx.Commit(ctx); err != nil {
			log.Println("Commit error")
            return 0, err
        }
        return 0, ErrRefreshReuse
	}

	if !expiresAt.After(now){
		return 0, ErrInvalidRefresh
	}

	const insQ = `
    INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
    VALUES ($1, $2, $3)
    RETURNING id
    `
    var newID int64
    if err := tx.QueryRow(ctx, insQ, userID, r.hashToken(newRaw), newExpiresAt).Scan(&newID); err != nil {
		log.Println("inserting error")
        return 0, err
    }
	
	const updQ = `
    UPDATE refresh_tokens
    SET revoked_at = $2,
        replaced_by = $3
    WHERE id = $1
    `
    if _, err := tx.Exec(ctx, updQ, oldID, now, newID); err != nil {
		log.Println("updating revoked_at and replace_by error")
        return 0, err
    }

    if err := tx.Commit(ctx); err != nil {
		log.Println("Commit error")
        return 0, err
    }
	

	return userID, nil
}

func (r *RefreshRepo) Revoke(ctx context.Context, rawToken string, now time.Time) error {
const q = `
	UPDATE refresh_tokens
	SET revoked_at = $2
	WHERE token_hash = $1
	  AND revoked_at IS NULL
	`
	_, err := r.pool.Exec(ctx, q, r.hashToken(rawToken), now)
	return err
}

func (r *RefreshRepo) RevokeAllForUser(ctx context.Context, userID int64, now time.Time) error {
	const q = `
    UPDATE refresh_tokens
    SET revoked_at = $2
    WHERE user_id = $1
      AND revoked_at IS NULL
    `
	_, err := r.pool.Exec(ctx, q, userID, now)
	return err
}