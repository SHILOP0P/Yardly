package auth

import (
	"net/http"
	"strings"
	"log"
	"context"

	"github.com/SHILOP0P/Yardly/backend/internal/httpx"
)

type AccessGuard interface {
    GetAccessState(ctx context.Context, userID int64) (tokenVersion int64, banned bool, err error)
}



func Middleware(jwtSvc *JWT, guard AccessGuard) func(http.Handler) http.Handler{
	return  func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := r.Header.Get("Authorization")
			
			if h == "" {
				log.Println("Authorization header is missing")
				httpx.WriteError(w, http.StatusUnauthorized, "missing authorization header")
				return
			}
			const prefix = "Bearer "
			if !strings.HasPrefix(h, prefix) {
				log.Println("Invalid authorization scheme")
				httpx.WriteError(w, http.StatusUnauthorized, "invalid authorization scheme")
				return
			}
			tokenStr := strings.TrimSpace(strings.TrimPrefix(h,prefix))
			userID, role, tv, err := jwtSvc.Parse(tokenStr)
			if err!=nil{
				log.Println("Invalid token:", err)
				httpx.WriteError(w, http.StatusUnauthorized, "invalid token")
				return 
			}

			dbTV, bannedNow, err :=guard.GetAccessState(r.Context(), userID)
			if err!=nil{
				 httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
                return
			}
			if bannedNow {
                httpx.WriteError(w, http.StatusForbidden, "banned")
                return
            }
            if dbTV != tv {
                httpx.WriteError(w, http.StatusUnauthorized, "token revoked")
                return
            }

			ctx:=WithUserID(r.Context(), userID)
			ctx=WithRole(ctx, role)
			ctx=WithBanned(ctx, bannedNow)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
