package pgrepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/SHILOP0P/Yardly/backend/internal/item"
)

type Repo struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

// GetByID: SELECT одной строки.
// Если строки нет — возвращаем item.ErrNotFound (а не pgx.ErrNoRows).
func (r *Repo) GetByID(ctx context.Context, id int64) (item.Item, error) {
	const q = `
SELECT id, owner_id, title, status, mode
FROM items
WHERE id = $1
`

	var it item.Item
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&it.ID,
		&it.OwnerID,
		&it.Title,
		&it.Status,
		&it.Mode,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return item.Item{}, item.ErrNotFound
		}
		return item.Item{}, fmt.Errorf("items pgrepo: get by id: %w", err)
	}

	return it, nil
}

func (r *Repo) List(ctx context.Context, f item.ListFilter) ([]item.Item, error) {
	limit := f.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	offset := f.Offset
	if offset < 0 {
		offset = 0
	}

	st := make([]string, 0, len(f.Status))
	for _, s:= range f.Status{
		st = append(st, string(s))
	}
	
	
	q := `
		SELECT id, owner_id, title, status, mode
		FROM items
		WHERE status = ANY($1::text[])
		`

		
	args := make([]any, 0, 4)
	n := 1

	args = append(args, st)
	n++

	if f.Mode != nil{
		q += fmt.Sprintf(" AND mode = $%d\n", n)
		args = append(args, *f.Mode)
		n++
	}
	
	q+= fmt.Sprintf(" ORDER BY id DESC LIMIT $%d OFFSET $%d", n, n+1)
	args = append(args, limit, offset)


	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("items pgrepo: list: %w", err)
	}
	defer rows.Close()

	var out []item.Item
	for rows.Next() {
		var it item.Item
		if err := rows.Scan(
			&it.ID,
			&it.OwnerID,
			&it.Title,
			&it.Status,
			&it.Mode,
		); err != nil {
			return nil, fmt.Errorf("items pgrepo: list scan: %w", err)
		}
		out = append(out, it)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("items pgrepo: list rows: %w", err)
	}

	return out, nil
}


func (r *Repo) Create(ctx context.Context, it *item.Item) error{
	const q = `
	INSERT INTO items (owner_id, title, status, mode)
	VALUES ($1, $2, $3, $4)
	RETURNING id
	`
	err := r.pool.QueryRow(ctx, q,
        it.OwnerID,
        it.Title,
        it.Status,
        it.Mode,
    ).Scan(&it.ID)
	if err != nil {
		return fmt.Errorf("items pgrepo: %w", err)
	}
	return nil
}

func (r *Repo) ListByOwnerPublic(ctx context.Context, ownerID int64, limit, offset int)([]item.Item, error){
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	const q = `
	SELECT id, owner_id, title, status, mode
	FROM items
	WHERE owner_id = $1
	AND status IN ('active', 'in_use')
	ORDER BY id DESC
	LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, q, ownerID, limit, offset)
	if err!=nil{
		return nil, fmt.Errorf("items pgrepo: list by owner public: %w", err)
	}
	defer rows.Close()

	out:=make([]item.Item, 0, limit)
	for rows.Next(){
		var it item.Item
		if err := rows.Scan(			
			&it.ID,
			&it.OwnerID,
			&it.Title,
			&it.Status,
			&it.Mode,);err!=nil{
			return nil, fmt.Errorf("items pgrepo: list by owner public scan: %w", err)
		}
		out = append(out, it)
	}
	if err:= rows.Err(); err!= nil{
		return nil, fmt.Errorf("items pgrepo: list by owner public rows: %w", err)
	}
	return out, nil
}


func (r *Repo) ListMyItems(ctx context.Context, ownerID int64, limit, offset int)([]item.Item, error){
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	const q = `
	SELECT id, owner_id, title, status, mode
	FROM items
	WHERE owner_id = $1
	AND status NOT IN ('deleted','transferred')
	ORDER BY id DESC
	LIMIT $2 OFFSET $3
	`
	rows, err := r.pool.Query(ctx, q, ownerID, limit, offset)
	if err!=nil{
		return nil, fmt.Errorf("items pgrepo: list my items: %w", err)
	}
	defer rows.Close()

	out := make([]item.Item, 0, 64)
	for rows.Next(){
		var it item.Item
		if err:=rows.Scan(
			&it.ID,
			&it.OwnerID,
			&it.Title,
			&it.Status,
			&it.Mode,
			); err!=nil{
			return nil, fmt.Errorf("items pgrepo: list my items scan: %w", err)
		}
		out = append(out, it)
	}

	if err := rows.Err(); err != nil{
		return nil, fmt.Errorf("items pgrepo: list my items rows: %w", err)
	}

	return out, nil
}