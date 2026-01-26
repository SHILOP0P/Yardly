package booking

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler

func RegisterRoutes(mux *http.ServeMux, repo Repo, items ItemGetter, authMw Middleware){//тут остановился
	h := NewHandler(repo, items)

	mux.Handle("POST /api/items/{id}/bookings", authMw(http.HandlerFunc(h.CreateRent)))
	mux.HandleFunc("GET /api/items/{id}/bookings", h.ListBusyForItem)

	mux.Handle("POST /api/bookings/{id}/approve", authMw(http.HandlerFunc(h.Approve)))
	mux.Handle("POST /api/bookings/{id}/return", authMw(http.HandlerFunc(h.Return)))
	mux.Handle("POST /api/bookings/{id}/handover", authMw(http.HandlerFunc(h.Handover)))

	mux.Handle("POST /api/bookings/{id}/cancel", authMw(http.HandlerFunc(h.Cancel)))


}