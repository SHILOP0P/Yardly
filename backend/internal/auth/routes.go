package auth

import (
	"net/http"
	"time"
)

type Mw func(http.Handler) http.Handler

func RegisterRoutes(mux *http.ServeMux, jwtSvc *JWT, refreshesRepo *RefreshRepo, refreshTTL time.Duration, users UserAuthenticator, authMw Mw){
	h := NewHandler(jwtSvc, refreshesRepo, refreshTTL, users)

	mux.HandleFunc("POST /api/auth/login", h.Login)

	mux.HandleFunc("POST /api/auth/refresh", h.Refresh)
	mux.Handle("POST /api/auth/logout", authMw(http.HandlerFunc(h.Logout)))
	mux.Handle("POST /api/auth/logout_all", authMw(http.HandlerFunc(h.LogoutAll)))
}