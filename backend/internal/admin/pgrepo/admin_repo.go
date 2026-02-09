package pgrepo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"


	"github.com/SHILOP0P/Yardly/backend/internal/admin"
)

type Repo struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repo { return &Repo{pool: pool} }

func (r *Repo) ListUsers(ctx context.Context, q string, limit, offset int) ([]admin.UserListItem, error) {
	// Поиск по email (ILIKE). q может быть пустым.
	const sqlQ = `
		SELECT
		id,
		email,
		role,
		banned_at,
		ban_reason,
		created_at,
		updated_at
		FROM users
		WHERE ($1 = '' OR email ILIKE '%' || $1 || '%')
		ORDER BY id DESC
		LIMIT $2 OFFSET $3
		`
	rows, err := r.pool.Query(ctx, sqlQ, q, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("admin list users: %w", err)
	}
	defer rows.Close()

	out := make([]admin.UserListItem, 0, limit)
	for rows.Next(){
		var u admin.UserListItem
		if err := rows.Scan(
			&u.ID,
			&u.Email,
			&u.Role,
			&u.BannedAt,
			&u.BanReason,
			&u.CreatedAt,
			&u.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("admin list users scan: %w", err)
		}
		out = append(out, u)
	}
	if err:= rows.Err(); err!=nil{
		return nil, fmt.Errorf("admin list users rows: %w", err)
	}

	return out, nil
}