package favorite

import (
	"net/http"
	_"github.com/SHILOP0P/Yardly/backend/internal/auth"
)

type Middleware func(http.Handler) http.Handler

func RegisterRoutes(mux *http.ServeMux, repo Repo, authMw Middleware){
	h := NewHandler(repo)

	mux.Handle("POST /api/items/{id}/favorite", authMw(http.HandlerFunc(h.Add)))
	mux.Handle("DELETE /api/items/{id}/favorite", authMw(http.HandlerFunc(h.Remove)))
	mux.Handle("GET /api/my/favorites", authMw(http.HandlerFunc(h.ListMy)))

	// опционально
	mux.Handle("GET /api/items/{id}/favorite", authMw(http.HandlerFunc(h.IsFavorite)))
}