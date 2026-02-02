package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"time"

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

func (r *RefreshRepo) Rotate(ctx context.Context, oldRaw, newRaw string, newExpiresAt, now time.Time)(int64, error){
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err!=nil{
		return 0, err
	}
	defer func(){
		_=tx.Rollback(ctx)
	}()

	const consumeQ = `
	UPDATE refresh_tokens
	SET revoked_at = $2
	WHERE token_hash = $1
	  AND revoked_at IS NULL
	  AND expires_at > $2
	RETURNING user_id
	`
	var userID int64
	err = tx.QueryRow(ctx, consumeQ, r.hashToken(oldRaw), now).Scan(&userID)
	if err != nil {
		return 0, err 
	}
	const insertQ = `
	INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
	VALUES ($1, $2, $3)
	`
	_, err = tx.Exec(ctx, insertQ, userID, r.hashToken(newRaw), newExpiresAt)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(ctx); err != nil {
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