package favorite

import(
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/SHILOP0P/Yardly/backend/internal/auth"
	"github.com/SHILOP0P/Yardly/backend/internal/httpx"
)

type Handler struct{
	repo Repo
}

func NewHandler(repo Repo) *Handler{return &Handler{repo: repo}}

func (h *Handler) Add(w http.ResponseWriter, r *http.Request){
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok{
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	itemID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err!=nil|| itemID <= 0 {
		httpx.WriteError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	f, err:= h.repo.Add(r.Context(), userID,itemID)
	if err!=nil{
	switch err {
		case ErrAlreadyExists:
			httpx.WriteError(w, http.StatusConflict, "already in favorites")
		case ErrNotFound:
			httpx.WriteError(w, http.StatusNotFound, "not found")
		default:
			httpx.WriteError(w, http.StatusInternalServerError, "internal error")
		}
		return
	}
	httpx.WriteJSON(w, http.StatusOK, f)
}

func (h *Handler) Remove(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	itemID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || itemID <= 0 {
		httpx.WriteError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	if err := h.repo.Remove(r.Context(), userID, itemID); err != nil {
		if err == ErrNotFound {
			httpx.WriteError(w, http.StatusNotFound, "not found")
			return
		}
		httpx.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListMy(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	limit := httpx.QueryInt(r, "limit", 20, 1, 100)
	offset := httpx.QueryInt(r, "offset", 0, 0, 1_000_000)

	items, err := h.repo.List(r.Context(), userID, limit, offset)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, items)
}

func (h *Handler) IsFavorite(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	itemID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || itemID <= 0 {
		httpx.WriteError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	okFav, err := h.repo.IsFavorite(r.Context(), userID, itemID)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	json.NewEncoder(w).Encode(map[string]bool{"is_favorite": okFav})
}