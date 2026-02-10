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


const selectAdminBookingCols = `
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

func scanAdminBooking(rs interface{ Scan(...any) error }, b *admin.AdminBooking) error {
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





func (r *Repo) ListUsers(ctx context.Context, q string, limit, offset int) ([]admin.UserListItem, error) {
	// Поиск по email (ILIKE). q может быть пустым.
	const sqlQ = `
		SELECT
		id,
		email,
		role,
		banned_at,
		ban_reason,
		ban_expires_at,
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
			&u.BanExpiresAt,
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


func (r *Repo) GetUser(ctx context.Context, id int64) (admin.UserListItem, error){
	const q = `
		SELECT
			id,
			email,
			role,
			banned_at,
			ban_reason,
			created_at,
			updated_at
		FROM users
		WHERE id = $1
	`
	var u admin.UserListItem
	if err := r.pool.QueryRow(ctx, q, id).Scan(
		&u.ID,
		&u.Email,
		&u.Role,
		&u.BannedAt,
		&u.BanReason,
		&u.CreatedAt,
		&u.UpdatedAt,
	); err != nil {
		return admin.UserListItem{}, fmt.Errorf("admin get user: %w", err)
	}
	return u, nil
}


func (r *Repo) PatchUser(ctx context.Context, actorAdminID, targetUserID int64, req admin.PatchUserRequest) (admin.UserListItem, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return admin.UserListItem{}, fmt.Errorf("admin patch user begin: %w", err)
	}
	defer tx.Rollback(ctx)

	// 1) select old FOR UPDATE
	const selQ = `
SELECT id, email, role, banned_at, ban_expires_at, ban_reason, created_at, updated_at
FROM users
WHERE id=$1
FOR UPDATE
`
	var old admin.UserListItem
	if err := tx.QueryRow(ctx, selQ, targetUserID).Scan(
		&old.ID, &old.Email, &old.Role, &old.BannedAt, &old.BanExpiresAt, &old.BanReason, &old.CreatedAt, &old.UpdatedAt,
	); err != nil {
		return admin.UserListItem{}, fmt.Errorf("admin patch user select: %w", err)
	}

	changed := make([]string, 0, 4)

	// 2) apply updates
	if req.Role != nil {
		const q = `UPDATE users SET role=$2, updated_at=now() WHERE id=$1`
		if _, err := tx.Exec(ctx, q, targetUserID, *req.Role); err != nil {
			return admin.UserListItem{}, fmt.Errorf("admin patch user role: %w", err)
		}
		changed = append(changed, "role")
	}

	if req.Ban != nil {
		if *req.Ban {
			const q = `
UPDATE users
SET banned_at=now(),
    ban_reason=$2,
    ban_expires_at=$3,
    updated_at=now()
WHERE id=$1
`
			if _, err := tx.Exec(ctx, q, targetUserID, req.BanReason, req.BanExpiresAt); err != nil {
				return admin.UserListItem{}, fmt.Errorf("admin patch user ban: %w", err)
			}
		} else {
			const q = `
UPDATE users
SET banned_at=NULL,
    ban_reason=NULL,
    ban_expires_at=NULL,
    updated_at=now()
WHERE id=$1
`
			if _, err := tx.Exec(ctx, q, targetUserID); err != nil {
				return admin.UserListItem{}, fmt.Errorf("admin patch user unban: %w", err)
			}
		}
		changed = append(changed, "banned_at", "ban_reason", "ban_expires_at")
	}

	// 3) read new
	const getQ = `
SELECT id, email, role, banned_at, ban_expires_at, ban_reason, created_at, updated_at
FROM users
WHERE id=$1
`
	var now admin.UserListItem
	if err := tx.QueryRow(ctx, getQ, targetUserID).Scan(
		&now.ID, &now.Email, &now.Role, &now.BannedAt, &now.BanExpiresAt, &now.BanReason, &now.CreatedAt, &now.UpdatedAt,
	); err != nil {
		return admin.UserListItem{}, fmt.Errorf("admin patch user reload: %w", err)
	}

	// 4) audit через CreateAdminEventTx
	if len(changed) > 0 {
		action := "user.patch"
		var reason *string
		if req.Ban != nil && *req.Ban && req.BanReason != nil {
			reason = req.BanReason
		}

		ev := admin.AdminEvent{
			ActorID:    actorAdminID,
			EntityType: "user",
			EntityID:   targetUserID,
			Action:     action,
			Reason:     reason,
			Meta: map[string]any{
				"changed_fields": changed,
				"old": map[string]any{
					"role":           old.Role,
					"banned_at":      old.BannedAt,
					"ban_expires_at": old.BanExpiresAt,
					"ban_reason":     old.BanReason,
				},
				"new": map[string]any{
					"role":           now.Role,
					"banned_at":      now.BannedAt,
					"ban_expires_at": now.BanExpiresAt,
					"ban_reason":     now.BanReason,
				},
			},
		}

		if err := r.CreateAdminEventTx(ctx, tx, ev); err != nil {
			return admin.UserListItem{}, fmt.Errorf("admin patch user audit: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return admin.UserListItem{}, fmt.Errorf("admin patch user commit: %w", err)
	}

	return now, nil
}

func(r *Repo) ListBookings(ctx context.Context, f admin.AdminBookingsFilter) ([]admin.AdminBooking, error){
		const q = `
	SELECT ` + selectAdminBookingCols + `
	FROM bookings
	WHERE
	($1::text IS NULL OR status = $1)
	AND ($2::text IS NULL OR type = $2)
	AND ($3::bigint IS NULL OR item_id = $3)
	AND ($4::bigint IS NULL OR requester_id = $4 OR owner_id = $4)
	ORDER BY created_at DESC, id DESC
	LIMIT $5 OFFSET $6
	`
	rows, err := r.pool.Query(ctx, q, f.Status, f.Type, f.ItemID, f.UserID, f.Limit, f.Offset)
	if err != nil {
		return nil, fmt.Errorf("admin list bookings: %w", err)
	}
	defer rows.Close()

	out := make([]admin.AdminBooking, 0, f.Limit)
	for rows.Next(){
		var b admin.AdminBooking
		if err:=scanAdminBooking(rows, &b); err != nil{
			return nil, fmt.Errorf("admin list bookings scan: %w", err)
		}
		out = append(out, b)
	}
	if err:=rows.Err();err!=nil{
		return nil, fmt.Errorf("admmin list bookings rows: %w", err)
	}
	return out, nil
}


func (r *Repo) GetBooking(ctx context.Context, id int64) (admin.AdminBooking, error) {
	const q = `
	SELECT ` + selectAdminBookingCols + `
	FROM bookings
	WHERE id = $1
	`
	var b admin.AdminBooking
	if err := scanAdminBooking(r.pool.QueryRow(ctx, q, id), &b); err != nil {
		return admin.AdminBooking{}, fmt.Errorf("admin get booking: %w", err)
	}
	return b, nil
}

