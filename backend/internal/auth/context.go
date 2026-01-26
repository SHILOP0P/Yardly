package auth

import "context"

type ctxKey int

const userIDKey ctxKey = 1

func WithUserID(ctx context.Context, userID int64) context.Context{
	return context.WithValue(ctx, userIDKey, userID)
}

func UserIDFromContext(ctx context.Context) (int64, bool){
	v := ctx.Value(userIDKey)
	id, ok := v.(int64)
	return id, ok
}