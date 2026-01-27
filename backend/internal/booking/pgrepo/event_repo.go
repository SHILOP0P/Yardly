package pgrepo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/SHILOP0P/Yardly/backend/internal/booking"
)

type EventRepo struct{}

func NewEventRepo() *EventRepo{return &EventRepo{}}


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