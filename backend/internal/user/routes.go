package user

import (
	"net/http"
	"github.com/SHILOP0P/Yardly/backend/internal/auth"
)


type Middleware func(http.Handler) http.Handler

func RegisterRoutes(mux *http.ServeMux, authMw Middleware, userRepo Repo, jwtSvc *auth.JWT){
	
	h := NewHandler(userRepo, jwtSvc)

	mux.HandleFunc("POST /api/auth/register", h.Register)
	mux.Handle("GET /api/users/me", authMw(http.HandlerFunc(h.Me)))
}