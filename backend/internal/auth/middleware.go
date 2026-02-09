package auth

import (
	"net/http"
	"strings"
	"log"

	"github.com/SHILOP0P/Yardly/backend/internal/httpx"
)

func Middleware(jwtSvc *JWT) func(http.Handler) http.Handler{
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
			userID, role, banned, err := jwtSvc.Parse(tokenStr)
			if err!=nil{
				log.Println("Invalid token:", err)
				httpx.WriteError(w, http.StatusUnauthorized, "invalid token")
				return 
			}
			ctx:=WithUserID(r.Context(), userID)
			ctx=WithRole(ctx, role)
			ctx=WithBanned(ctx, banned)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}