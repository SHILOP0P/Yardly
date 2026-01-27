package pgrepo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/SHILOP0P/Yardly/backend/internal/booking"
)

type Repo struct {
	pool *pgxpool.Pool
	eventRepo *EventRepo
}

func New(pool *pgxpool.Pool, eventRepo *EventRepo) *Repo {
    return &Repo{
        pool:      pool,
        eventRepo: eventRepo,
    }
}


const selectBookingCols = `
	id, item_id, requester_id, owner_id,
	type, status,
	start_at, end_at,
	handover_deadline,
	handover_confirmed_by_owner_at,
	handover_confirmed_by_requester_at,
	return_confirmed_by_owner_at,
	return_confirmed_by_requester_at,
	created_at
`

type rowScanrer interface{
	Scan(dest ...any) error
}

func scanBooking(rs rowScanrer, b *booking.Booking)error{
	return rs.Scan(
		&b.ID,
		&b.ItemID,
		&b.RequesterID,
		&b.OwnerID,
		&b.Type,
		&b.Status,
		&b.Start,
		&b.End,
		&b.HandoverDeadline,
		&b.HandoverConfirmedByOwnerAt,
		&b.HandoverConfirmedByRequesterAt,
		&b.ReturnConfirmedByOwnerAt,
		&b.ReturnConfirmedByRequesterAt,
		&b.CreatedAt,
	)
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
	q := `
SELECT ` + selectBookingCols + `
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
    &b.HandoverConfirmedByOwnerAt,
    &b.HandoverConfirmedByRequesterAt,
    &b.ReturnConfirmedByOwnerAt,
    &b.ReturnConfirmedByRequesterAt,
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
SELECT ` + selectBookingCols + `
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
		&b.HandoverConfirmedByOwnerAt,
		&b.HandoverConfirmedByRequesterAt,
		&b.ReturnConfirmedByOwnerAt,
		&b.ReturnConfirmedByRequesterAt,
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
	SELECT ` + selectBookingCols + `
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
		&b.HandoverConfirmedByOwnerAt,
		&b.HandoverConfirmedByRequesterAt,
		&b.ReturnConfirmedByOwnerAt,
		&b.ReturnConfirmedByRequesterAt,
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
  RETURNING id
`
	rows, err := tx.Query(ctx, declineQ,
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
	defer rows.Close()

	declinedIDs:=make([]int64, 0, 8)
	for rows.Next(){
		var id int64
		if err:= rows.Scan(&id); err!=nil{
			return booking.Booking{}, fmt.Errorf("bookings pgrepo: decline competitors scan: %w", err)
		}
		declinedIDs = append(declinedIDs, id)
	}
	if err:=rows.Err();err!=nil{
		return booking.Booking{}, fmt.Errorf("bookings pgrepo: decline competitors rows: %w", err)
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

	actor := ownerID
	fromApproved:=booking.StatusRequested
	toApproved := booking.StatusApproved

	if err :=r.eventRepo.InsertBookingEvent(ctx, tx, b.ID, &actor, "approve", &fromApproved, &toApproved, nil); err != nil {
		return booking.Booking{}, err
	}

	fromDecl := booking.StatusRequested
	toDecl := booking.StatusDeclined

	for _, id := range declinedIDs {
		// meta можно не делать, но полезно
		if err := r.eventRepo.InsertBookingEvent(ctx, tx, id, &actor, "auto_decline_competitor", &fromDecl, &toDecl, nil); err != nil {
			return booking.Booking{}, err
		}
	}


	if err :=tx.Commit(ctx); err!=nil{
		return booking.Booking{}, fmt.Errorf("bookings pgrepo: commit: %w", err)
	}
	return b, nil

}

func (r *Repo) ReturnRent(ctx context.Context, bookingID int64, actorID int64, now time.Time)(booking.Booking, error){
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return booking.Booking{}, fmt.Errorf("bookings pgrepo: begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	const selectQ = `
	SELECT ` + selectBookingCols + `
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
		&b.HandoverConfirmedByOwnerAt,
		&b.HandoverConfirmedByRequesterAt,
		&b.ReturnConfirmedByOwnerAt,
		&b.ReturnConfirmedByRequesterAt,
		&b.CreatedAt,
	)
	if err != nil{
		if errors.Is(err, pgx.ErrNoRows){
			return booking.Booking{}, booking.ErrNotFound
		}
		return booking.Booking{}, fmt.Errorf("bookings pgrepo: return select: %w", err)
	}
	if b.Type != booking.TypeRent{
		return booking.Booking{}, booking.ErrInvalidState
	}
	if actorID!=b.OwnerID&&actorID!=b.RequesterID{
		return booking.Booking{}, booking.ErrForbidden
	}
	if actorID==b.OwnerID&&actorID==b.RequesterID{
		return booking.Booking{}, booking.ErrInvalidState
	}
	if b.Status!=booking.StatusInUse&&b.Status!=booking.StatusReturnPending{
		return booking.Booking{}, booking.ErrInvalidState
	}
	if actorID == b.OwnerID && b.ReturnConfirmedByOwnerAt == nil {
		b.ReturnConfirmedByOwnerAt = &now
	}
	if actorID == b.RequesterID && b.ReturnConfirmedByRequesterAt == nil {
		b.ReturnConfirmedByRequesterAt = &now
	}

	const updMarks = `
	UPDATE bookings
	SET return_confirmed_by_owner_at = $1,
		return_confirmed_by_requester_at = $2
	WHERE id = $3
	RETURNING
		return_confirmed_by_owner_at,
		return_confirmed_by_requester_at,
		status
	`
	err = tx.QueryRow(ctx, updMarks,
		b.ReturnConfirmedByOwnerAt,
		b.ReturnConfirmedByRequesterAt,
		b.ID,
	).Scan(&b.ReturnConfirmedByOwnerAt, &b.ReturnConfirmedByRequesterAt, &b.Status)
	if err != nil {
		return booking.Booking{}, fmt.Errorf("bookings pgrepo: return update marks: %w", err)
	}

	actor := actorID

	by := "requester"
	if actorID == b.OwnerID {
		by = "owner"
	}
	meta, _ := json.Marshal(map[string]string{"by": by})

	if err := r.eventRepo.InsertBookingEvent(
		ctx, tx,
		b.ID,
		&actor,
		"return_confirm",
		nil,
		nil,
		meta,
	); err != nil {
		return booking.Booking{}, err
	}


	if b.ReturnConfirmedByOwnerAt != nil && b.ReturnConfirmedByRequesterAt != nil {
		from := b.Status
				
		const updStatus = `
		UPDATE bookings
		SET status = $1
		WHERE id = $2 AND status IN ($3, $4)
		RETURNING status
		`
		err = tx.QueryRow(ctx, updStatus,
			booking.StatusCompleted,
			b.ID,
			booking.StatusInUse,
			booking.StatusReturnPending,
		).Scan(&b.Status)
		if err != nil {
			return booking.Booking{}, fmt.Errorf("bookings pgrepo: return set status: %w", err)
		}

		to := b.Status

		if err := r.eventRepo.InsertBookingEvent(
			ctx, tx,
			b.ID,
			&actor,
			"status_change",
			&from,
			&to,
			nil,
		); err != nil {
			return booking.Booking{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return booking.Booking{}, fmt.Errorf("bookings pgrepo: commit: %w", err)
	}
	return b, nil	
}

func (r *Repo) HandoverRent(ctx context.Context, bookingID int64, actorID int64, now time.Time) (booking.Booking, error) {
	
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return booking.Booking{}, fmt.Errorf("bookings pgrepo: begin tx: %w", err)
	}
	defer tx.Rollback(ctx)


	// Если дедлайн прошёл — НЕ переводим в in_use.
	// Это защищает инвариант "только в течение 24 часов".
	const selectQ = `
	SELECT ` + selectBookingCols + `
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
		&b.HandoverConfirmedByOwnerAt,
		&b.HandoverConfirmedByRequesterAt,
		&b.ReturnConfirmedByOwnerAt,
		&b.ReturnConfirmedByRequesterAt,
		&b.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return booking.Booking{}, r.explainNoRows(ctx, bookingID, actorID, booking.StatusApproved)
		}
		return booking.Booking{}, fmt.Errorf("bookings pgrepo: handover: %w", err)
	}
	if b.Type!=booking.TypeRent{
		return booking.Booking{}, booking.ErrInvalidState
	}
	if actorID!=b.OwnerID&&actorID!=b.RequesterID{
		return booking.Booking{}, booking.ErrForbidden
	}
	if actorID==b.OwnerID&&actorID==b.RequesterID{
		return booking.Booking{}, booking.ErrInvalidState
	}
	if b.Status!= booking.StatusApproved{
		return booking.Booking{}, booking.ErrInvalidState
	}
	if b.HandoverDeadline == nil || now.After(*b.HandoverDeadline){
		return booking.Booking{}, booking.ErrInvalidState
	}
	if actorID == b.OwnerID&&b.HandoverConfirmedByOwnerAt == nil{
		b.HandoverConfirmedByOwnerAt = &now
	}
	if actorID == b.RequesterID && b.HandoverConfirmedByRequesterAt == nil {
		b.HandoverConfirmedByRequesterAt = &now
	}

	const updMarks = `
		UPDATE bookings
		SET handover_confirmed_by_owner_at = $1,
			handover_confirmed_by_requester_at = $2
		WHERE id = $3
		RETURNING
			handover_confirmed_by_owner_at,
			handover_confirmed_by_requester_at,
			status
		`
	err = tx.QueryRow(ctx, updMarks,
		b.HandoverConfirmedByOwnerAt,
		b.HandoverConfirmedByRequesterAt,
		b.ID,
	).Scan(&b.HandoverConfirmedByOwnerAt, &b.HandoverConfirmedByRequesterAt, &b.Status)
	if err != nil {
		return booking.Booking{}, fmt.Errorf("bookings pgrepo: handover update marks: %w", err)
	}

	actor:= actorID
	by := "requester"
	if actorID == b.OwnerID{
		by = "owner"
	}
	meta, _:=json.Marshal(map[string]string{"by":by})
	
	if err :=r.eventRepo.InsertBookingEvent(ctx, tx,
		b.ID,
		&actor,
		"handover_confirm",
		nil,
		nil,
		meta,
	); err != nil {
		return booking.Booking{}, err
	}

	if b.HandoverConfirmedByOwnerAt != nil && b.HandoverConfirmedByRequesterAt != nil {
		from := b.Status

		const updStatus = `
			UPDATE bookings
			SET status = $1
			WHERE id = $2 AND status = $3
			RETURNING status
			`
		err = tx.QueryRow(ctx, updStatus,
			booking.StatusInUse,
			b.ID,
			booking.StatusApproved,
		).Scan(&b.Status)
		if err != nil {
			return booking.Booking{}, fmt.Errorf("bookings pgrepo: handover set status: %w", err)
		}

		to := b.Status
		if err :=r.eventRepo.InsertBookingEvent(ctx, tx,
			b.ID,
			&actor,
			"status_change",
			&from,
			&to,
			nil,
		); err != nil {
			return booking.Booking{}, err
		}
	}
	if err:=tx.Commit(ctx); err!=nil{
		return booking.Booking{}, fmt.Errorf("bookings pgrepo: commit: %w", err)
	}

	return b, nil
}

func (r *Repo) ExpireOverdueHandovers(ctx context.Context, now time.Time) (int64, error) {
	tx, err:= r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err!=nil{
		return 0, fmt.Errorf("bookings pgrepo: begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	q := `
	UPDATE bookings
	SET status = $1
	WHERE type = $2
	AND status = $3
	AND handover_deadline IS NOT NULL
	AND handover_deadline < $4
	RETURNING id
	`
	rows, err := tx.Query(ctx, q, booking.StatusExpired, booking.TypeRent, booking.StatusApproved, now)
	if err != nil {
		return 0, fmt.Errorf("bookings pgrepo: expire query: %w", err)
	}
	defer rows.Close()

	var n int64
	from := booking.StatusApproved
	to:=booking.StatusExpired

	for rows.Next(){
		var id int64
		if err := rows.Scan(&id);err!=nil{
			return 0, fmt.Errorf("bookings pgrepo: expire scan: %w", err)
		}
		n++
		if err:=r.eventRepo.InsertBookingEvent(ctx, tx, id, nil, "expire", &from, &to, nil); err!=nil{
			return 0, err
		}
	}
	if err:=rows.Err(); err!=nil{
		return 0, fmt.Errorf("bookings pgrepo: expire rows: %w", err)
	}
	
	if err:=tx.Commit(ctx); err!=nil{
		return 0, fmt.Errorf("bookings pgrepo: commit: %w", err)
	}
	return n, nil
}

func (r *Repo) CancelRent(ctx context.Context, bookingID, requesterID int64)(booking.Booking, error){
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return booking.Booking{}, fmt.Errorf("bookings pgrepo: begin tx: %w", err)
	}
	defer tx.Rollback(ctx)
	const q = `
	UPDATE bookings
	SET status = $3
	WHERE id = $1 AND requester_id = $2 AND status = $4 AND type = $5
	RETURNING `+ selectBookingCols + `
	`
	var out booking.Booking
	err=scanBooking(tx.QueryRow(ctx,q,
	bookingID,
		requesterID,
		booking.StatusCanceled,
		booking.StatusRequested,
		booking.TypeRent,
	), &out)
	if err != nil{
		if errors.Is(err, pgx.ErrNoRows){
			return booking.Booking{}, r.explainCancelNoRows(ctx, bookingID, requesterID)
		}
		return booking.Booking{}, err
	}
	actor := requesterID
	from := booking.StatusRequested
	to := booking.StatusCanceled

	if err := r.eventRepo.InsertBookingEvent(ctx, tx, out.ID, &actor, "cancel", &from, &to, nil); err != nil {
		return booking.Booking{}, err
	}

	if err:=tx.Commit(ctx); err != nil{
		return booking.Booking{}, fmt.Errorf("bookings pgrepo: commit: %w", err)
	}
	return out, nil
}

func(r *Repo) ListMyBookings(ctx context.Context, requesterID int64, statuses[]booking.Status, limit, offset int)([]booking.Booking, error){
	if limit<=0{
		limit = 20
	}
	if limit>100{
		limit = 100
	}
	if offset <0{
		offset = 0
	}
	st := make([]string, 0, len(statuses))
	for _, s:= range statuses{
		st = append(st, string(s))
	}

	const base = `
	SELECT ` + selectBookingCols + `
	FROM bookings
	WHERE requester_id = $1
	`
	var rows pgx.Rows
	var err error

	if len(st) == 0{
		const q = base+`ORDER BY created_at DESC, id DESC
		LIMIT $2 OFFSET $3
		`
		rows, err = r.pool.Query(ctx, q, requesterID, limit, offset)
	} else{
		const q = base + `
			AND status = ANY($2::text[])
			ORDER BY created_at DESC, id DESC
			LIMIT $3 OFFSET $4
			`
		rows, err = r.pool.Query(ctx, q, requesterID, st, limit, offset)
	}
	if err != nil{
		return nil, fmt.Errorf("bookings pgrepo: list my bookings: %w", err)
	}
	defer rows.Close()

	out := make([]booking.Booking, 0, limit)
	for rows.Next(){
		var b booking.Booking
		if err:= scanBooking(rows, &b); err != nil {
			return nil, fmt.Errorf("bookings pgrepo: list my bookings scan: %w", err)
		}
		out = append(out, b)
	}
	if err := rows.Err(); err != nil{
		return nil, fmt.Errorf("bookings pgrepo: list my bookings rows: %w", err)
	}
	return out, nil
}


func(r *Repo) ListMyItemsBookings(ctx context.Context, ownerID int64, statuses[]booking.Status, limit, offset int)([]booking.Booking, error){
	if limit<=0{
		limit = 20
	}
	if limit>100{
		limit = 100
	}
	if offset <0{
		offset = 0
	}
	st := make([]string, 0, len(statuses))
	for _, s:= range statuses{
		st = append(st, string(s))
	}

	const base = `
	SELECT ` + selectBookingCols + `
	FROM bookings
	WHERE owner_id = $1
	`
	var rows pgx.Rows
	var err error

	if len(st) == 0{
		const q = base+`ORDER BY created_at DESC, id DESC
		LIMIT $2 OFFSET $3
		`
		rows, err = r.pool.Query(ctx, q, ownerID, limit, offset)
	} else{
		const q = base + `
			AND status = ANY($2::text[])
			ORDER BY created_at DESC, id DESC
			LIMIT $3 OFFSET $4
			`
		rows, err = r.pool.Query(ctx, q, ownerID, st, limit, offset)
	}
	if err != nil{
		return nil, fmt.Errorf("bookings pgrepo: list my bookings: %w", err)
	}
	defer rows.Close()

	out := make([]booking.Booking, 0, limit)
	for rows.Next(){
		var b booking.Booking
		if err:=scanBooking(rows, &b); err != nil {
			return nil, fmt.Errorf("bookings pgrepo: list my bookings scan: %w", err)
		}
		out = append(out, b)
	}
	if err := rows.Err(); err != nil{
		return nil, fmt.Errorf("bookings pgrepo: list my bookings rows: %w", err)
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