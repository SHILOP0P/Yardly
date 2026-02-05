package pgrepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/SHILOP0P/Yardly/backend/internal/favorite"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct{
	pool *pgxpool.Pool
}
func New(pool *pgxpool.Pool) *Repo{
	return &Repo{pool: pool}
}

func (r *Repo) Add(ctx context.Context,  userID, itemID int64)(favorite.Favorite, error){
	const q = `
	INSERT INTO favorites(user_id, item_id)
	VALUES ($1, $2)
	RETURNING user_id, item_id, created_at
	`
	var f favorite.Favorite
	err := r.pool.QueryRow(ctx, q, userID, itemID).Scan(&f.UserID, &f.ItemID, &f.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return favorite.Favorite{}, favorite.ErrAlreadyExists
		}
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			// FK violation: user или item не существует
			return favorite.Favorite{}, favorite.ErrNotFound
		}
		return favorite.Favorite{}, fmt.Errorf("favorites add: %w", err)
	}
	return f, nil
}

func (r *Repo) Remove(ctx context.Context, userID, itemID int64)error{
	const q = `DELETE FROM favorites WHERE user_id=$1 AND item_id=$2`
	ct, err := r.pool.Exec(ctx, q, userID, itemID)
	if err != nil {
		return fmt.Errorf("favorites remove: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return favorite.ErrNotFound
	}
	return nil
}

func (r *Repo) IsFavorite(ctx context.Context, userID, itemID int64)(bool, error){
	const q = `SELECT 1 FROM favorites WHERE user_id=$1 AND item_id=$2`
	var one int
	err := r.pool.QueryRow(ctx, q, userID, itemID).Scan(&one)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	return false, fmt.Errorf("favorites isFavorite: %w", err)
}

func (r *Repo) List(ctx context.Context, userID int64, limit, offset int)([]favorite.FavoriteItem, error){
	const q = `
	SELECT
	  f.item_id,
	  i.title,
	  i.status,
	  i.mode,
	  i.owner_id,
	  f.created_at
	FROM favorites f
	JOIN items i ON i.id = f.item_id
	WHERE f.user_id = $1
	  AND i.status IN ('active','in_use') -- скрываем archived/deleted/transferred
	ORDER BY f.created_at DESC
	LIMIT $2 OFFSET $3
	`
	rows, err := r.pool.Query(ctx, q, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("favorites list: %w", err)
	}
	defer rows.Close()

	out := make([]favorite.FavoriteItem, 0, limit)
	for rows.Next(){
		var x favorite.FavoriteItem
		if err:=rows.Scan(&x.ItemID, &x.Title, &x.Status, &x.Mode, &x.OwnerID, &x.FavoritedAt);err!=nil{
			return nil, fmt.Errorf("favorites list scan: %w", err)
		}
		out = append(out, x)
	}
	if err:=rows.Err();err!=nil{
		return nil, fmt.Errorf("favorites list rows: %w", err)
	}
	return out, nil
}