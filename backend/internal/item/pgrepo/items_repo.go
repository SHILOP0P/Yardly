package pgrepo

import (
	"context"
	"errors"
	"fmt"
	"strings"

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
SELECT id, owner_id, title, status, mode, description, price, deposit, location, category
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
		&it.Description,
		&it.Price,
		&it.Deposit,
		&it.Location,
		&it.Category,
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
	for _, s := range f.Status {
		st = append(st, string(s))
	}

	q := `
		SELECT id, owner_id, title, status, mode, description, price, deposit, location, category
		FROM items
		WHERE status = ANY($1::text[])
		`

	args := make([]any, 0, 8)
	n := 1

	args = append(args, st)
	n++

	if f.Mode != nil {
		q += fmt.Sprintf(" AND mode = $%d\n", n)
		args = append(args, *f.Mode)
		n++
	}

	if f.Category != nil {
		q += fmt.Sprintf(" AND category ILIKE $%d\n", n)
		args = append(args, "%"+*f.Category+"%")
		n++
	}

	if f.Location != nil {
		q += fmt.Sprintf(" AND location ILIKE $%d\n", n)
		args = append(args, "%"+*f.Location+"%")
		n++
	}

	if f.MinPrice != nil {
		q += fmt.Sprintf(" AND price >= $%d\n", n)
		args = append(args, *f.MinPrice)
		n++
	}

	if f.MaxPrice != nil {
		q += fmt.Sprintf(" AND price <= $%d\n", n)
		args = append(args, *f.MaxPrice)
		n++
	}

	q += fmt.Sprintf(" ORDER BY id DESC LIMIT $%d OFFSET $%d", n, n+1)
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
			&it.Description,
			&it.Price,
			&it.Deposit,
			&it.Location,
			&it.Category,
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

func (r *Repo) Create(ctx context.Context, it *item.Item) error {
	const q = `
		INSERT INTO items (
			owner_id,
			title,
			status,
			mode,
			description,
			price,
			deposit,
			location,
			category
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		RETURNING id
	`

	err := r.pool.QueryRow(ctx, q,
		it.OwnerID,
		it.Title,
		it.Status,
		it.Mode,
		it.Description,
		it.Price,
		it.Deposit,
		it.Location,
		it.Category,
	).Scan(&it.ID)
	if err != nil {
		return fmt.Errorf("items pgrepo create: %w", err)
	}

	if len(it.Images) > 0 {
		if err := r.insertImages(ctx, it.ID, it.Images); err != nil {
			return err
		}
	}

	return nil
}

func (r *Repo) insertImages(ctx context.Context, itemID int64, imgs []item.ItemImage) error {
	const q = `
		INSERT INTO item_images (item_id, url, sort_order)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	for i := range imgs {
		imgs[i].ItemID = itemID
		if imgs[i].SortOrder == 0 {
			imgs[i].SortOrder = i + 1
		}

		if err := r.pool.QueryRow(ctx, q, imgs[i].ItemID, imgs[i].URL, imgs[i].SortOrder).
			Scan(&imgs[i].ID, &imgs[i].CreatedAt); err != nil {
			return fmt.Errorf("items pgrepo insertImages: %w", err)
		}
	}

	return nil
}

func (r *Repo) ListByOwnerPublic(ctx context.Context, ownerID int64, f item.ListFilter) ([]item.Item, error) {
	limit := f.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	offset := f.Offset
	if offset < 0 {
		offset = 0
	}

	q := `
	SELECT id, owner_id, title, status, mode, description, price, deposit, location, category
	FROM items
	WHERE owner_id = $1
	AND status IN ('active', 'in_use')
	`

	args := make([]any, 0, 7)
	n := 2
	args = append(args, ownerID)

	if f.Mode != nil {
		q += fmt.Sprintf(" AND mode = $%d\n", n)
		args = append(args, *f.Mode)
		n++
	}

	if f.Category != nil {
		q += fmt.Sprintf(" AND category ILIKE $%d\n", n)
		args = append(args, "%"+*f.Category+"%")
		n++
	}

	if f.Location != nil {
		q += fmt.Sprintf(" AND location ILIKE $%d\n", n)
		args = append(args, "%"+*f.Location+"%")
		n++
	}

	if f.MinPrice != nil {
		q += fmt.Sprintf(" AND price >= $%d\n", n)
		args = append(args, *f.MinPrice)
		n++
	}

	if f.MaxPrice != nil {
		q += fmt.Sprintf(" AND price <= $%d\n", n)
		args = append(args, *f.MaxPrice)
		n++
	}

	q += fmt.Sprintf(" ORDER BY id DESC LIMIT $%d OFFSET $%d", n, n+1)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("items pgrepo: list by owner public: %w", err)
	}
	defer rows.Close()

	out := make([]item.Item, 0, limit)
	for rows.Next() {
		var it item.Item
		if err := rows.Scan(
			&it.ID,
			&it.OwnerID,
			&it.Title,
			&it.Status,
			&it.Mode,
			&it.Description,
			&it.Price,
			&it.Deposit,
			&it.Location,
			&it.Category,
		); err != nil {
			return nil, fmt.Errorf("items pgrepo: list by owner public scan: %w", err)
		}
		out = append(out, it)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("items pgrepo: list by owner public rows: %w", err)
	}
	return out, nil
}

func (r *Repo) ListMyItems(ctx context.Context, ownerID int64, f item.ListFilter) ([]item.Item, error) {
	limit := f.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	offset := f.Offset
	if offset < 0 {
		offset = 0
	}

	q := `
	SELECT id, owner_id, title, status, mode, description, price, deposit, location, category
	FROM items
	WHERE owner_id = $1
	AND status NOT IN ('deleted','transferred')
	`

	args := make([]any, 0, 7)
	n := 2
	args = append(args, ownerID)

	if f.Mode != nil {
		q += fmt.Sprintf(" AND mode = $%d\n", n)
		args = append(args, *f.Mode)
		n++
	}

	if f.Category != nil {
		q += fmt.Sprintf(" AND category ILIKE $%d\n", n)
		args = append(args, "%"+*f.Category+"%")
		n++
	}

	if f.Location != nil {
		q += fmt.Sprintf(" AND location ILIKE $%d\n", n)
		args = append(args, "%"+*f.Location+"%")
		n++
	}

	if f.MinPrice != nil {
		q += fmt.Sprintf(" AND price >= $%d\n", n)
		args = append(args, *f.MinPrice)
		n++
	}

	if f.MaxPrice != nil {
		q += fmt.Sprintf(" AND price <= $%d\n", n)
		args = append(args, *f.MaxPrice)
		n++
	}

	q += fmt.Sprintf(" ORDER BY id DESC LIMIT $%d OFFSET $%d", n, n+1)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("items pgrepo: list my items: %w", err)
	}
	defer rows.Close()

	out := make([]item.Item, 0, 64)
	for rows.Next() {
		var it item.Item
		if err := rows.Scan(
			&it.ID,
			&it.OwnerID,
			&it.Title,
			&it.Status,
			&it.Mode,
			&it.Description,
			&it.Price,
			&it.Deposit,
			&it.Location,
			&it.Category,
		); err != nil {
			return nil, fmt.Errorf("items pgrepo: list my items scan: %w", err)
		}
		out = append(out, it)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("items pgrepo: list my items rows: %w", err)
	}

	return out, nil
}

//	Images

func (r *Repo) ListImages(ctx context.Context, itemID int64) ([]item.ItemImage, error) {
	const q = `
	SELECT id, item_id, url, sort_order, created_at
	FROM item_images
	WHERE item_id = $1
	ORDER BY sort_order ASC
	`
	rows, err := r.pool.Query(ctx, q, itemID)
	if err != nil {
		return nil, fmt.Errorf("items pgrepo: list images: %w", err)
	}
	defer rows.Close()

	out := make([]item.ItemImage, 0, 8)
	for rows.Next() {
		var im item.ItemImage
		if err := rows.Scan(&im.ID, &im.ItemID, &im.URL, &im.SortOrder, &im.CreatedAt); err != nil {
			return nil, fmt.Errorf("items pgrepo: list images scan: %w", err)
		}
		out = append(out, im)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("items pgrepo: list images rows: %w", err)
	}
	return out, nil
}

func (r *Repo) AddImage(ctx context.Context, itemID int64, url string) (item.ItemImage, error) {
	url = strings.TrimSpace(url)
	if url == "" {
		return item.ItemImage{}, fmt.Errorf("items pgrepo: add image: empty url")
	}

	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return item.ItemImage{}, fmt.Errorf("items pgrepo: add image begin: %w", err)
	}
	defer tx.Rollback(ctx)

	const nextQ = `
SELECT sort_order
FROM item_images
WHERE item_id = $1
ORDER BY sort_order DESC
LIMIT 1
FOR UPDATE
`
	next := 1
	var last int
	if err := tx.QueryRow(ctx, nextQ, itemID).Scan(&last); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return item.ItemImage{}, fmt.Errorf("items pgrepo: add image next: %w", err)
		}
	} else {
		next = last + 1
	}

	const insQ = `
	INSERT INTO item_images (item_id, url, sort_order)
	VALUES ($1, $2, $3)
	RETURNING id, item_id, url, sort_order, created_at
	`
	var im item.ItemImage
	if err := tx.QueryRow(ctx, insQ, itemID, url, next).Scan(
		&im.ID, &im.ItemID, &im.URL, &im.SortOrder, &im.CreatedAt,
	); err != nil {
		return item.ItemImage{}, fmt.Errorf("items pgrepo: add image insert: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return item.ItemImage{}, fmt.Errorf("items pgrepo: add image commit: %w", err)
	}
	return im, nil
}

func (r *Repo) DeleteImage(ctx context.Context, itemID, imageID int64) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("items pgrepo: delete image begin: %w", err)
	}
	defer tx.Rollback(ctx)

	const selQ = `
SELECT sort_order
FROM item_images
WHERE id = $1 AND item_id = $2
FOR UPDATE
`
	var ord int
	if err := tx.QueryRow(ctx, selQ, imageID, itemID).Scan(&ord); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return item.ErrNotFound
		}
		return fmt.Errorf("items pgrepo: delete image select: %w", err)
	}

	const delQ = `DELETE FROM item_images WHERE id = $1 AND item_id = $2`
	ct, err := tx.Exec(ctx, delQ, imageID, itemID)
	if err != nil {
		return fmt.Errorf("items pgrepo: delete image delete: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return item.ErrNotFound
	}

	const shiftQ = `
UPDATE item_images
SET sort_order = sort_order - 1
WHERE item_id = $1 AND sort_order > $2
`
	if _, err := tx.Exec(ctx, shiftQ, itemID, ord); err != nil {
		return fmt.Errorf("items pgrepo: delete image shift: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("items pgrepo: delete image commit: %w", err)
	}
	return nil
}
