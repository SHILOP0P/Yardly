package auth

import (
	"net/http"
	"time"
)

func RegisterRoutes(mux *http.ServeMux, jwtSvc *JWT, refreshesRepo *RefreshRepo, refreshTTL time.Duration, users UserAuthenticator,){
	h := NewHandler(jwtSvc, refreshesRepo, refreshTTL, users)

	mux.HandleFunc("POST /api/auth/login", h.Login)

	mux.HandleFunc("POST /api/auth/refresh", h.Refresh)
	mux.HandleFunc("POST /api/auth/logout", h.Logout)

}