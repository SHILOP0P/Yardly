package auth

import "net/http"

func RegisterRoutes(mux *http.ServeMux, jwtSvc *JWT){
	h := NewHandler(jwtSvc)

	mux.HandleFunc("POST /api/auth/login", h.Login)
}