package admin

import(
	"net/http"
	
)

type adminHandler func(http.Handler) http.Handler

func RegisterRoutes(mux *http.ServeMux, adminRepo Repo, adminChain adminHandler){
	h := New(adminRepo)
	//User
	mux.Handle("GET /api/admin/users", adminChain(http.HandlerFunc(h.ListUsers)))
	mux.Handle("GET /api/admin/users/{id}", adminChain(http.HandlerFunc(h.GetUser)))
	mux.Handle("PATCH /api/admin/users/{id}", adminChain(http.HandlerFunc(h.PatchUser)))

	//Booking
	mux.Handle("GET /api/admin/bookings", adminChain(http.HandlerFunc(h.ListBookings)))
	mux.Handle("GET /api/admin/bookings/{id}", adminChain(http.HandlerFunc(h.GetBooking)))
	mux.Handle("GET /api/admin/bookings/{id}/events", adminChain(http.HandlerFunc(h.ListBookingEvents)))

	//events
	mux.Handle("GET /api/admin/events", adminChain(http.HandlerFunc(h.ListAdminEvents)))
}