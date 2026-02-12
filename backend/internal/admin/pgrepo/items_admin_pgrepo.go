package pgrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/SHILOP0P/Yardly/backend/internal/admin"
	"github.com/jackc/pgx/v5"
)

const selectAdminItemCols = `
id, owner_id, title, status, mode, blocked_at, block_reason
`

type rowScanner interface { Scan(...any) error}

func scanAdminItem(rs rowScanner, it *admin.AdminItem) error{
	return rs.Scan(
		&it.ID,
		&it.OwnerID,
		&it.Title,
		&it.Status,
		&it.Mode,
		&it.BlockedAt,
		&it.BlockReason,
	)
}

func (r *Repo) GetItem(ctx context.Context, id int64)(admin.AdminItem, error){
	const q = `SELECT ` + selectAdminItemCols + ` FROM items WHERE id=$1`
	var it admin.AdminItem
	if err := scanAdminItem(r.pool.QueryRow(ctx, q, id), &it); err != nil {
		return admin.AdminItem{}, fmt.Errorf("admin get item: %w", err)
	}
	return it, nil
}

func (r *Repo) ListItems(ctx context.Context, f admin.AdminItemsFilter)([]admin.AdminItem, error){
	if f.Limit <= 0 {
		f.Limit = 50
	}
	if f.Limit > 200 {
		f.Limit = 200
	}
	if f.Offset < 0 {
		f.Offset = 0
	}

	args := make([]any, 0, 0)
	n:=1

	q := `SELECT ` + selectAdminItemCols + ` FROM items WHERE 1=1`

	if f.Q !=nil && *f.Q!=""{
		q += fmt.Sprintf(" AND title ILIKE '%%' || $%d || '%%'", n)
		args = append(args, *f.Q)
		n++
	}
	if f.Status!=nil&&*f.Status!=""{
		q += fmt.Sprintf(" AND status = $%d", n)
		args = append(args, *f.Status)
		n++
	}
	if f.Mode != nil && *f.Mode != "" {
		q += fmt.Sprintf(" AND mode = $%d", n)
		args = append(args, *f.Mode)
		n++
	}
	if f.IncludeArchived||f.IncludeDeleted||f.IncludeTransferred{
		allowed:=make([]string, 0, 5)
		allowed = append(allowed, "active", "in_use")
		if f.IncludeArchived {
		allowed = append(allowed, "archived")
		}
		if f.IncludeDeleted {
			allowed = append(allowed, "deleted")
		}
		if f.IncludeTransferred {
			allowed = append(allowed, "transferred")
		}
		q+=fmt.Sprintf(" AND status = ANY($%d::text[])", n)
		args = append(args, allowed)
		n++
	}

	q += fmt.Sprintf(" ORDER BY id DESC LIMIT $%d OFFSET $%d", n, n+1)
	args = append(args, f.Limit, f.Offset)

	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("admin list items: %w", err)
	}
	defer rows.Close()

	out := make([]admin.AdminItem, 0, f.Limit)
	for rows.Next() {
		var it admin.AdminItem
		if err := scanAdminItem(rows, &it); err != nil {
			return nil, fmt.Errorf("admin list items scan: %w", err)
		}
		out = append(out, it)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("admin list items rows: %w", err)
	}
	return out, nil
	
}

func (r *Repo) PatchItem(ctx context.Context, actorAdminID, itemID int64, req admin.PatchItemRequest)(admin.AdminItem, error){
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err!=nil{
		return admin.AdminItem{}, fmt.Errorf("admin patch item begin: %w", err)
	}
	defer tx.Rollback(ctx)

	const sel = `SELECT ` + selectAdminItemCols + ` FROM items WHERE id=$1 FOR UPDATE`
	var old admin.AdminItem
	if err := scanAdminItem(tx.QueryRow(ctx, sel, itemID), &old); err != nil {
		return admin.AdminItem{}, fmt.Errorf("admin patch item select: %w", err)
	}

	changed := make([]string, 0, 3)
	if req.Title != nil {
		const uq = `UPDATE items SET title=$2 WHERE id=$1`
		if _, err := tx.Exec(ctx, uq, itemID, *req.Title); err != nil {
			return admin.AdminItem{}, fmt.Errorf("admin patch item title: %w", err)
		}
		changed = append(changed, "title")
	}

	// mode
	// ТУТ МОЖЕТ ПЕРЕДАВАТЬСЯ НЕСКОЛЬКО СТАТУСОВ
	if req.Mode != nil {
		// валидируем по твоим mode
		switch *req.Mode {
		case "sale", "rent", "free", "sale_rent":
		default:
			return admin.AdminItem{}, fmt.Errorf("invalid mode")
		}
		const uq = `UPDATE items SET mode=$2 WHERE id=$1`
		if _, err := tx.Exec(ctx, uq, itemID, *req.Mode); err != nil {
			return admin.AdminItem{}, fmt.Errorf("admin patch item mode: %w", err)
		}
		changed = append(changed, "mode")
	}

	// status
	if req.Status != nil {
		switch *req.Status {
		case "active", "in_use", "archived", "transferred":
		default:
			return admin.AdminItem{}, fmt.Errorf("invalid status")
		}
		const uq = `UPDATE items SET status=$2 WHERE id=$1`
		if _, err := tx.Exec(ctx, uq, itemID, *req.Status); err != nil {
			return admin.AdminItem{}, fmt.Errorf("admin patch item status: %w", err)
		}
		changed = append(changed, "status")
	}

	const get = `SELECT ` + selectAdminItemCols + ` FROM items WHERE id=$1`
	var now admin.AdminItem
	if err := scanAdminItem(tx.QueryRow(ctx, get, itemID), &now); err != nil {
		return admin.AdminItem{}, fmt.Errorf("admin patch item reload: %w", err)
	}

	// audit
	if len(changed) > 0 {
		ev := admin.AdminEvent{
			ActorID:    actorAdminID,
			EntityType: "item",
			EntityID:   itemID,
			Action:     "item.patch",
			Reason:     nil,
			Meta: map[string]any{
				"changed_fields": changed,
				"old": old,
				"new": now,
			},
		}
		if err := r.CreateAdminEventTx(ctx, tx, ev); err != nil {
			return admin.AdminItem{}, fmt.Errorf("admin patch item audit: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return admin.AdminItem{}, fmt.Errorf("admin patch item commit: %w", err)
	}
	return now, nil
}

func (r *Repo) BlockItem(ctx context.Context, actorAdminID, itemID int64, reason *string) (admin.AdminItem, error){
	now := time.Now().UTC()

	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err!=nil{
		return admin.AdminItem{}, fmt.Errorf("admin patch item begin: %w", err)
	}
	defer tx.Rollback(ctx)

	const sel =`SELECT ` + selectAdminItemCols + ` FROM items WHERE id=$1 FOR UPDATE`
	var old admin.AdminItem
	if err := scanAdminItem(tx.QueryRow(ctx, sel, itemID), &old); err != nil {
		return admin.AdminItem{}, fmt.Errorf("admin block item select: %w", err)
	}

	const upd = `UPDATE items SET blocked_at=$2, block_reason=$3 WHERE id=$1`
	if _, err := tx.Exec(ctx, upd, itemID, now, reason); err != nil {
		return admin.AdminItem{}, fmt.Errorf("admin block item update: %w", err)
	}

	const get = `SELECT ` + selectAdminItemCols + ` FROM items WHERE id=$1`
	var cur admin.AdminItem
	if err := scanAdminItem(tx.QueryRow(ctx, get, itemID), &cur); err != nil {
		return admin.AdminItem{}, fmt.Errorf("admin block item reload: %w", err)
	}

	ev := admin.AdminEvent{
		ActorID:    actorAdminID,
		EntityType: "item",
		EntityID:   itemID,
		Action:     "item.block",
		Reason:     reason,
		Meta: map[string]any{
			"old": old,
			"new": cur,
		},
	}
	if err := r.CreateAdminEventTx(ctx, tx, ev); err != nil {
		return admin.AdminItem{}, fmt.Errorf("admin block item audit: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return admin.AdminItem{}, fmt.Errorf("admin block item commit: %w", err)
	}
	return cur, nil
}

func (r *Repo) UnblockItem(ctx context.Context, actorAdminID, itemID int64, reason *string) (admin.AdminItem, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return admin.AdminItem{}, fmt.Errorf("admin unblock item begin: %w", err)
	}
	defer tx.Rollback(ctx)

	const sel = `SELECT ` + selectAdminItemCols + ` FROM items WHERE id=$1 FOR UPDATE`
	var old admin.AdminItem
	if err := scanAdminItem(tx.QueryRow(ctx, sel, itemID), &old); err != nil {
		return admin.AdminItem{}, fmt.Errorf("admin unblock item select: %w", err)
	}

	const upd = `UPDATE items SET blocked_at=NULL, block_reason=NULL WHERE id=$1`
	if _, err := tx.Exec(ctx, upd, itemID); err != nil {
		return admin.AdminItem{}, fmt.Errorf("admin unblock item update: %w", err)
	}

	const get = `SELECT ` + selectAdminItemCols + ` FROM items WHERE id=$1`
	var cur admin.AdminItem
	if err := scanAdminItem(tx.QueryRow(ctx, get, itemID), &cur); err != nil {
		return admin.AdminItem{}, fmt.Errorf("admin unblock item reload: %w", err)
	}

	ev := admin.AdminEvent{
		ActorID:    actorAdminID,
		EntityType: "item",
		EntityID:   itemID,
		Action:     "item.unblock",
		Reason:     reason,
		Meta: map[string]any{
			"old": old,
			"new": cur,
		},
	}
	if err := r.CreateAdminEventTx(ctx, tx, ev); err != nil {
		return admin.AdminItem{}, fmt.Errorf("admin unblock item audit: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return admin.AdminItem{}, fmt.Errorf("admin unblock item commit: %w", err)
	}
	return cur, nil
}

func (r *Repo) DeleteItem(ctx context.Context, actorAdminID, itemID int64, reason *string)(admin.AdminItem, error){
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return admin.AdminItem{}, fmt.Errorf("admin delete item begin: %w", err)
	}
	defer tx.Rollback(ctx)

	const sel = `SELECT ` + selectAdminItemCols + ` FROM items WHERE id=$1 FOR UPDATE`
	var old admin.AdminItem
	if err := scanAdminItem(tx.QueryRow(ctx, sel, itemID), &old); err != nil {
		return admin.AdminItem{}, fmt.Errorf("admin delete item select: %w", err)
	}

	const upd = `UPDATE items SET status='deleted' WHERE id=$1`
	if _, err := tx.Exec(ctx, upd, itemID); err != nil {
		return admin.AdminItem{}, fmt.Errorf("admin delete item update: %w", err)
	}

	const get = `SELECT ` + selectAdminItemCols + ` FROM items WHERE id=$1`
	var cur admin.AdminItem
	if err := scanAdminItem(tx.QueryRow(ctx, get, itemID), &cur); err != nil {
		return admin.AdminItem{}, fmt.Errorf("admin delete item reload: %w", err)
	}

	ev := admin.AdminEvent{
		ActorID:    actorAdminID,
		EntityType: "item",
		EntityID:   itemID,
		Action:     "item.delete",
		Reason:     reason,
		Meta: map[string]any{
			"old": old,
			"new": cur,
		},
	}

	if err := r.CreateAdminEventTx(ctx, tx, ev); err != nil {
		return admin.AdminItem{}, fmt.Errorf("admin delete item audit: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return admin.AdminItem{}, fmt.Errorf("admin delete item commit: %w", err)
	}
	return cur, nil
}