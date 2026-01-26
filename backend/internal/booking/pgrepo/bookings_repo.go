package pgrepo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/SHILOP0P/Yardly/backend/internal/booking"
)

type Repo struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

func (r *Repo) Create(ctx context.Context, b *booking.Booking) error {
	// created_at ставится DEFAULT now() в таблице
	const q = `
INSERT INTO bookings (
	item_id, requester_id, owner_id,
	type, status,
	start_at, end_at,
	handover_deadline
) VALUES (
	$1, $2, $3,
	$4, $5,
	$6, $7,
	$8
)
RETURNING id, created_at
`
	err := r.pool.QueryRow(ctx, q,
		b.ItemID,
		b.RequesterID,
		b.OwnerID,
		b.Type,
		b.Status,
		b.Start,
		b.End,
		b.HandoverDeadline,
	).Scan(&b.ID, &b.CreatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505"{
			return booking.ErrDuplicateActiveRequest
		}
		return err
	}
	return nil
}

func (r *Repo) GetByID(ctx context.Context, id int64) (booking.Booking, error) {
	const q = `
SELECT
	id, item_id, requester_id, owner_id,
	type, status,
	start_at, end_at,
	handover_deadline,
	created_at
FROM bookings
WHERE id = $1
`
	var b booking.Booking
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&b.ID,
		&b.ItemID,
		&b.RequesterID,
		&b.OwnerID,
		&b.Type,
		&b.Status,
		&b.Start,
		&b.End,
		&b.HandoverDeadline,
		&b.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return booking.Booking{}, booking.ErrNotFound
		}
		return booking.Booking{}, fmt.Errorf("bookings pgrepo: get by id: %w", err)
	}
	return b, nil
}

func (r *Repo) ListByItem(ctx context.Context, itemID int64) ([]booking.Booking, error) {
	const q = `
SELECT
	id, item_id, requester_id, owner_id,
	type, status,
	start_at, end_at,
	handover_deadline,
	created_at
FROM bookings
WHERE item_id = $1
ORDER BY id DESC
`
	rows, err := r.pool.Query(ctx, q, itemID)
	if err != nil {
		return nil, fmt.Errorf("bookings pgrepo: list by item: %w", err)
	}
	defer rows.Close()

	out := make([]booking.Booking, 0, 16)
	for rows.Next() {
		var b booking.Booking
		if err := rows.Scan(
			&b.ID,
			&b.ItemID,
			&b.RequesterID,
			&b.OwnerID,
			&b.Type,
			&b.Status,
			&b.Start,
			&b.End,
			&b.HandoverDeadline,
			&b.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("bookings pgrepo: list by item scan: %w", err)
		}
		out = append(out, b)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("bookings pgrepo: list by item rows: %w", err)
	}

	return out, nil
}

func (r *Repo) ApproveRent(ctx context.Context, bookingID int64, ownerID int64)(booking.Booking, error){
	tx, err := r.pool.BeginTx(ctx,pgx.TxOptions{})
	if err != nil{
		return booking.Booking{}, fmt.Errorf("bookings pgrepo: begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	const selectQ = `
	SELECT
		id, item_id, requester_id, owner_id,
		type, status,
		start_at, end_at,
		handover_deadline,
		created_at
	FROM bookings
	WHERE id = $1
	FOR UPDATE
	`
	var b booking.Booking
	err = tx.QueryRow(ctx, selectQ, bookingID).Scan(
		&b.ID,
		&b.ItemID,
		&b.RequesterID,
		&b.OwnerID,
		&b.Type,
		&b.Status,
		&b.Start,
		&b.End,
		&b.HandoverDeadline,
		&b.CreatedAt,
	)
	if err != nil{
		if errors.Is(err,pgx.ErrNoRows){
			return  booking.Booking{}, booking.ErrNotFound
		}
		return booking.Booking{}, fmt.Errorf("bookings pgrepo: approve select: %w", err)
	}
	if b.OwnerID != ownerID {
		return booking.Booking{}, booking.ErrForbidden
	}
	if b.Type != booking.TypeRent{
		return booking.Booking{}, booking.ErrInvalidState
	}
	if b.Status != booking.StatusRequested {
		return booking.Booking{}, booking.ErrInvalidState
	}
	if b.Start == nil || b.End == nil {
		return booking.Booking{}, fmt.Errorf("rent booking must have start/end")
	}

	dedline := b.Start.Add(24*time.Hour)

	const declineQ = `
UPDATE bookings
SET status = $1
WHERE item_id = $2
  AND type = $3
  AND status = $4
  AND id <> $5
  AND start_at < $6
  AND end_at   > $7
`
	_, err = tx.Exec(ctx, declineQ,
		booking.StatusDeclined,
		b.ItemID,
		booking.TypeRent,
		booking.StatusRequested,
		b.ID,
		*b.End,
		*b.Start,
	)
	if err != nil{
		return  booking.Booking{}, fmt.Errorf("bookings pgrepo: decline competitors: %w", err)
	}

	const approveQ = `
	UPDATE bookings
	SET status = $1,
		handover_deadline = $2
	WHERE id = $3
	RETURNING status, handover_deadline
	`
	err = tx.QueryRow(ctx, approveQ, booking.StatusApproved, dedline, b.ID).Scan(&b.Status,&b.HandoverDeadline)
	if err!=nil{
		if errors.Is(err, pgx.ErrNoRows) {
			// теоретически редкий кейс, но пусть будет
			return booking.Booking{}, booking.ErrNotFound
		}
		return booking.Booking{}, err
	}
	if err :=tx.Commit(ctx); err!=nil{
		return booking.Booking{}, fmt.Errorf("bookings pgrepo: commit: %w", err)
	}
	return b, nil

}

func (r *Repo) ReturnRent(ctx context.Context, bookingID int64, ownerID int64)(booking.Booking, error){
	const q = `
UPDATE bookings
SET status = $1
WHERE id = $2
  AND owner_id = $3
  AND type = $4
  AND status IN ($5, $6)
RETURNING
	id, item_id, requester_id, owner_id,
	type, status,
	start_at, end_at,
	handover_deadline,
	created_at
`
	var b booking.Booking
	err := r.pool.QueryRow(ctx, q,
		booking.StatusCompleted,
		bookingID,
		ownerID,
		booking.TypeRent,
		booking.StatusInUse,
		booking.StatusReturnPending,
	).Scan(
		&b.ID, &b.ItemID, &b.RequesterID, &b.OwnerID,
		&b.Type, &b.Status,
		&b.Start, &b.End,
		&b.HandoverDeadline,
		&b.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows){
			return booking.Booking{}, r.explainNoRows(ctx, bookingID, ownerID, booking.StatusInUse)
		}
		return booking.Booking{}, fmt.Errorf("bookings pgrepo: return: %w", err)
	}
	return b, nil
}

func (r *Repo) HandoverRent(ctx context.Context, bookingID int64, ownerID int64, now time.Time) (booking.Booking, error) {
	// Если дедлайн прошёл — НЕ переводим в in_use.
	// Это защищает инвариант "только в течение 24 часов".
	const q = `
UPDATE bookings
SET status = $1
WHERE id = $2
  AND owner_id = $3
  AND type = $4
  AND status = $5
  AND handover_deadline IS NOT NULL
  AND $6 <= handover_deadline
RETURNING
	id, item_id, requester_id, owner_id,
	type, status,
	start_at, end_at,
	handover_deadline,
	created_at
`
	var b booking.Booking
	err := r.pool.QueryRow(ctx, q,
		booking.StatusInUse,
		bookingID,
		ownerID,
		booking.TypeRent,
		booking.StatusApproved,
		now,
	).Scan(
		&b.ID, &b.ItemID, &b.RequesterID, &b.OwnerID,
		&b.Type, &b.Status,
		&b.Start, &b.End,
		&b.HandoverDeadline,
		&b.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return booking.Booking{}, r.explainNoRows(ctx, bookingID, ownerID, booking.StatusApproved)
		}
		return booking.Booking{}, fmt.Errorf("bookings pgrepo: handover: %w", err)
	}

	return b, nil
}

func (r *Repo) ExpireOverdueHandovers(ctx context.Context, now time.Time) (int64, error) {
	const q = `
UPDATE bookings
SET status = $1
WHERE type = $2
  AND status = $3
  AND handover_deadline IS NOT NULL
  AND handover_deadline < $4
`
	ct, err := r.pool.Exec(ctx, q,
		booking.StatusExpired,
		booking.TypeRent,
		booking.StatusApproved,
		now,
	)
	if err != nil {
		return 0, fmt.Errorf("bookings pgrepo: expire overdue: %w", err)
	}
	return ct.RowsAffected(), nil
}

func (r *Repo) CancelRent(ctx context.Context, bookingID, requesterID int64)(booking.Booking, error){
	const q = `
	UPDATE bookings
	SET status = $3
	WHERE id = $1 AND requester_id = $2 AND status = $4 AND type = $5
	RETURNING id, item_id, requester_id, owner_id, type, status, start_at, end_at, handover_deadline, created_at
	`
	var out booking.Booking
	err:=r.pool.QueryRow(ctx,q,
	bookingID,
		requesterID,
		booking.StatusCanceled,
		booking.StatusRequested,
		booking.TypeRent,
	).Scan(
		&out.ID, &out.ItemID, &out.RequesterID, &out.OwnerID, &out.Type, &out.Status,
		&out.Start, &out.End, &out.HandoverDeadline, &out.CreatedAt,
	)
	if err != nil{
		if errors.Is(err, pgx.ErrNoRows){
			return booking.Booking{}, r.explainCancelNoRows(ctx, bookingID, requesterID)
		}
		return booking.Booking{}, err
	}
	return out, nil
}




















func (r *Repo) explainNoRows(ctx context.Context, bookingID, ownerID int64, expectedStatus booking.Status) error{
	const q = `
	SELECT owner_id, status
	FROM bookings
	WHERE id = $1
	`

	var dbOwnerID int64
	var dbStatus booking.Status

	err := r.pool.QueryRow(ctx, q, bookingID).Scan(&dbOwnerID, &dbStatus)
	if err!= nil{
		if errors.Is(err, pgx.ErrNoRows){
			return booking.ErrNotFound
		}
		return err
	}
	if dbOwnerID!=ownerID{
		return booking.ErrForbidden
	}
	if dbStatus!= expectedStatus{
		return  booking.ErrInvalidState
	}
	return booking.ErrInvalidState
}

func (r *Repo) explainCancelNoRows(ctx context.Context, bookingID, requesterID int64) error{
	const q = `
		SELECT requester_id, status, type
		FROM bookings
		WHERE id = $1
		`
	var dbRequesterID int64
	var dbStatus booking.Status
	var dbType booking.Type

	err:=r.pool.QueryRow(ctx, q,bookingID).Scan(&dbRequesterID, &dbStatus, &dbType)
	if err != nil{
		if errors.Is(err, pgx.ErrNoRows){
			return booking.ErrNotFound
		}
		return err
	}
	if dbRequesterID != requesterID {
		return booking.ErrForbidden
	}
	// если это не rent или статус не requested — это конфликт состояния
	return booking.ErrInvalidState
}