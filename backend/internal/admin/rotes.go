package admin

import(
	"net/http"

)

type adminHandler func(http.Handler) http.Handler

func RegisterRoutes(mux *http.ServeMux, adminRepo Repo, adminChain adminHandler){
	h := New(adminRepo)

	mux.Handle("GET /api/admin/users", adminChain(http.HandlerFunc(h.ListUsers)))
}