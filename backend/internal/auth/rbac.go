package auth

import (
	"net/http"

	"github.com/SHILOP0P/Yardly/backend/internal/httpx"
)


func RequireAdmin(next http.Handler) http.Handler {
	return requireRole(next, RoleAdmin)
}

func RequireSuperAdmin(next http.Handler) http.Handler {
	return requireRole(next, RoleSuperAdmin)
}

func requireRole(next http.Handler, min Role) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// userID должен быть уже в контексте после Middleware(jwtSvc)
		if _, ok := UserIDFromContext(r.Context()); !ok {
			httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		if BannedFromContext(r.Context()){
			httpx.WriteError(w, http.StatusForbidden, "banned")
			return
		}

		role, ok := RoleFromContext(r.Context())
		if !ok {
			httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		if !roleAtLeast(role, min) {
			httpx.WriteError(w, http.StatusForbidden, "forbidden")
			return
		}

		next.ServeHTTP(w, r)
	})
}


func roleAtLeast(got, min Role) bool{
	rank := func(r Role) int{
		switch r {
		case RoleUser:
			return 1
		case RoleAdmin:
			return 2
		case RoleSuperAdmin:
			return 3
		default:
			return 0
		}

	}
	return rank(got)>=rank(min)
}