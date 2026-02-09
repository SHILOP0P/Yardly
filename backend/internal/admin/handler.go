package admin

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/SHILOP0P/Yardly/backend/internal/httpx"
)

type Handler struct {
	repo Repo
}

func New(repo Repo) *Handler{
	return &Handler{repo: repo}
}

// GET /api/admin/users?q=&limit=&offset=
func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request){
	q := r.URL.Query().Get("q")

	limit := 50
	if s := r.URL.Query().Get("limit"); s != "" {
		if v, err := strconv.Atoi(s); err == nil {
			limit = v
		}
	}
	offset := 0
	if s := r.URL.Query().Get("offset"); s != "" {
		if v, err := strconv.Atoi(s); err == nil {
			offset = v
		}
	}
	if limit <= 0 || limit > 200 {
		httpx.WriteError(w, http.StatusBadRequest, "invalid limit")
		return
	}
	if offset < 0 {
		httpx.WriteError(w, http.StatusBadRequest, "invalid offset")
		return
	}

	users, err := h.repo.ListUsers(r.Context(), q, limit, offset)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]any{
		"users":  users,
		"limit":  limit,
		"offset": offset,
		"q":      q,
	})
}