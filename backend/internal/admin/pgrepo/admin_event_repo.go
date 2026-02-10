package pgrepo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5"
	bookingpg "github.com/SHILOP0P/Yardly/backend/internal/booking/pgrepo"
	"github.com/SHILOP0P/Yardly/backend/internal/admin"
	"github.com/SHILOP0P/Yardly/backend/internal/booking"
)

type execer interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	QueryRow(context.Context, string, ...any) pgx.Row
}

func (r *Repo) CreateAdminEventTx(ctx context.Context, q execer, ev admin.AdminEvent) error {
	const ins = `
INSERT INTO admin_events (actor_user_id, entity_type, entity_id, action, reason, meta)
VALUES ($1, $2, $3, $4, $5, $6)
`
	var metaBytes []byte
	if ev.Meta != nil {
		b, err := json.Marshal(ev.Meta)
		if err != nil {
			return fmt.Errorf("admin event marshal meta: %w", err)
		}
		metaBytes = b
	}

	if _, err := q.Exec(ctx, ins,
		ev.ActorID,
		ev.EntityType,
		ev.EntityID,
		ev.Action,
		ev.Reason,
		metaBytes,
	); err != nil {
		return fmt.Errorf("admin event insert: %w", err)
	}

	return nil
}


func (r *Repo) ListAdminEvents(ctx context.Context, f admin.AdminEventsFilter)([]admin.AdminEvent, error){
	const q = `
SELECT
  id,
  actor_user_id,
  entity_type,
  entity_id,
  action,
  reason,
  meta,
  created_at
FROM admin_events
WHERE
  ($1::text  IS NULL OR entity_type = $1)
  AND ($2::bigint IS NULL OR entity_id = $2)
  AND ($3::bigint IS NULL OR actor_user_id = $3)
ORDER BY id DESC
LIMIT $4 OFFSET $5
`
	rows, err := r.pool.Query(ctx, q, f.EntityType, f.EntityID, f.ActorUserID, f.Limit, f.Offset)
	if err != nil {
		return nil, fmt.Errorf("admin list events: %w", err)
	}
	defer rows.Close()

	out := make([]admin.AdminEvent, 0, f.Limit)
	for rows.Next(){
		var e admin.AdminEvent
		var metaBytes []byte
		if err := rows.Scan(&e.ID,
			&e.ActorID,
			&e.EntityType,
			&e.EntityID,
			&e.Action,
			&e.Reason,
			&metaBytes,
			&e.CreatedAt,); err!=nil{
				return nil, fmt.Errorf("admin list events scan: %w", err)
		}
		if len(metaBytes)>0{
			var v any
			if err := json.Unmarshal(metaBytes, &v); err != nil{
				e.Meta = string(metaBytes)
			} else{
				e.Meta = v
			}
		}
		out = append(out, e)
	}

	if err := rows.Err(); err!= nil{
		return nil, fmt.Errorf("admin list events rows: %w", err)
	}
	return out, nil
}

func (r *Repo) ListBookingEvents(ctx context.Context, bookingID int64, limit, offset int) ([]booking.Event, error) {
	er := bookingpg.NewEventRepo() // важно: это booking/pgrepo EventRepo
	return er.ListBookingEvents(ctx, r.pool, bookingID, limit, offset)
}

