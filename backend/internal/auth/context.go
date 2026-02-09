package auth

import "context"

type ctxKey int

const (
	userIDKey ctxKey = iota + 1
	roleKey
	bannedKey
)

func WithUserID(ctx context.Context, userID int64) context.Context{
	return context.WithValue(ctx, userIDKey, userID)
}

func UserIDFromContext(ctx context.Context) (int64, bool){
	v := ctx.Value(userIDKey)
	id, ok := v.(int64)
	return id, ok
}

func WithRole(ctx context.Context, role Role) context.Context{
	return context.WithValue(ctx, roleKey, role)
}

func RoleFromContext(ctx context.Context) (Role, bool) {
	v:=ctx.Value(roleKey)
	r, ok:=v.(Role)
	return r, ok
}


func WithBanned(ctx context.Context, banned bool) context.Context {
	return context.WithValue(ctx, bannedKey, banned)
}

func BannedFromContext(ctx context.Context) bool {
	v := ctx.Value(bannedKey)
	b, ok := v.(bool)
	if !ok {
		return false
	}
	return b
}