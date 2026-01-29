package pgrepo

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/SHILOP0P/Yardly/backend/internal/booking"
)

type EventRepo struct{}

func NewEventRepo() *EventRepo{return &EventRepo{}}


const selectEventCols = `
SELECT
  id, booking_id, actor_user_id,
  action, from_status, to_status,
  meta, created_at
FROM booking_events
`


func (r *EventRepo) InsertBookingEvent(ctx context.Context, tx pgx.Tx, bookingID int64, actorID *int64, action string, from *booking.Status, to *booking.Status, meta []byte)error{
	const q = `
	INSERT INTO booking_events (booking_id, actor_user_id, action, from_status, to_status, meta)
	VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := tx.Exec(ctx, q, bookingID,
		actorID, action, from, to, meta,)
	
	if err != nil{
		return fmt.Errorf("bookings pgrepo: insert booking event: %w", err)
	}
	return nil
}

func (r *EventRepo) ListBookingEvents(ctx context.Context, pool interface{
	Query(context.Context, string, ... any)(pgx.Rows, error)}, bookingID int64, limit, offset int)([]booking.Event, error){
		if limit <= 0{
			limit = 50
		}
		if limit >200{
			limit = 200
		}
		if offset < 0{
			offset = 0
		}

		const q = `
		`+selectEventCols+`
		WHERE booking_id = $1
		ORDER BY created_at ASC, id ASC
		LIMIT $2 OFFSET $3
		`
		rows, err := pool.Query(ctx, q, bookingID, limit, offset)
		if err != nil{
			return nil, fmt.Errorf("event repo: list booking events: %w", err)
		}

		defer rows.Close()
		out := make([]booking.Event, 0, limit)
		for rows.Next(){
			var e booking.Event
			var metaBytes []byte

			if err:= rows.Scan(&e.ID,
			&e.BookingID,
			&e.ActorUserID,
			&e.Action,
			&e.FromStatus,
			&e.ToStatus,
			&metaBytes,
			&e.CreatedAt,); err != nil{
				return nil, fmt.Errorf("event repo: scan booking event: %w", err)
			}
			if len(metaBytes)>0{
				var m map[string]any
				if err := json.Unmarshal(metaBytes, &m); err == nil{
					e.Meta = m
				} else {
					e.Meta = string(metaBytes)
				}
			}
			out = append(out, e)
		}
		if err = rows.Err(); err != nil{
			return nil, fmt.Errorf("event repo: rows booking events: %w", err)
		}
		return out, nil
}