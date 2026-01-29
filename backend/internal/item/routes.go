package item

import "net/http"

type Middleware func(http.Handler) http.Handler

func RegisterRoutes(mux *http.ServeMux, repo Repo, authMw Middleware) {
	h := NewHandler(repo)

	mux.Handle("POST /api/items", authMw(http.HandlerFunc(h.Create)))
	mux.HandleFunc("GET /api/items", h.List)
	mux.HandleFunc("GET /api/items/{id}", h.GetByID)

	mux.Handle("GET /api/my/items", authMw(http.HandlerFunc(h.ListMyItems)))
	mux.HandleFunc("GET /api/users/{id}/items", h.ListByOwnerPublic)

}
