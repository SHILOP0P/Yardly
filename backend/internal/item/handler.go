package item

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"log"

	"github.com/SHILOP0P/Yardly/backend/internal/auth"
	"github.com/SHILOP0P/Yardly/backend/internal/httpx"
)

type Handler struct {
	repo Repo
}

 
func NewHandler(repo Repo) *Handler {
	return &Handler{repo: repo}
}



func (h *Handler) List (w http.ResponseWriter, r *http.Request){
	q := r.URL.Query()

	var f ListFilter

	f.Status = []Status{StatusActive, StatusInUse}

	if v := q.Get("mode"); v != "" {
		m := DealMode(v)
		if !m.Valid() {
			httpx.WriteError(w, http.StatusBadRequest, "invalid mode")
			return
		}
		f.Mode = &m
	}

	if v := q.Get("limit"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			httpx.WriteError(w, http.StatusBadRequest, "invalid limit")
			return
		}
		f.Limit = n
	}

	// offset
	if v := q.Get("offset"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			httpx.WriteError(w, http.StatusBadRequest, "invalid offset")
			return
		}
		f.Offset = n
	}

	items, err := h.repo.List(r.Context(), f)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.WriteJSON(w, http.StatusOK, items)
}


func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	it, err := h.repo.GetByID(r.Context(), id)
	if err != nil{
		if errors.Is(err, ErrNotFound){
			httpx.WriteError(w, http.StatusNotFound, "item not found")
			return
		}
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.WriteJSON(w, http.StatusOK, it)
}

func (h *Handler)Create(w http.ResponseWriter, r *http.Request){
	ownerID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
        return
	}

	var dto struct {
		Title  string   `json:"title"`
        Status Status   `json:"status"`
        Mode   DealMode `json:"mode"`
	}
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil{
		httpx.WriteError(w, http.StatusBadRequest, "invalid json body")
        return
	}
	if dto.Title == "" || !dto.Status.Valid() || !dto.Mode.Valid(){
		httpx.WriteError(w, http.StatusBadRequest, "invalid item fields")
		return
	}

	it := Item{
		 OwnerID: ownerID,
        Title:   dto.Title,
        Status:  dto.Status,
        Mode:    dto.Mode,
	}
	if err := h.repo.Create(r.Context(), &it); err !=nil{
		log.Println("item create error:", err)
    	httpx.WriteError(w, http.StatusInternalServerError, "could not create item")
        return
	}
	httpx.WriteJSON(w, http.StatusCreated, it)
}


func (h *Handler) ListMyItems(w http.ResponseWriter, r *http.Request){
	ownerID, ok := auth.UserIDFromContext(r.Context())
	if !ok{
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	limit, offset, err := parseLimitOffset(r)
	if err!=nil{
		httpx.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	items, err := h.repo.ListMyItems(r.Context(), ownerID, limit, offset)
	if err!=nil{
		log.Println("list my items error:", err)
		httpx.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, items)
}

func (h *Handler) ListByOwnerPublic(w http.ResponseWriter, r *http.Request){
	ownerID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || ownerID <= 0{
		httpx.WriteError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	limit, offset, err := parseLimitOffset(r)
	if err!=nil{
		httpx.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	items, err := h.repo.ListByOwnerPublic(r.Context(), ownerID, limit, offset)
	if err!=nil{
		log.Println("list owner items error:", err)
		httpx.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}


	httpx.WriteJSON(w, http.StatusOK, items)
}






func parseLimitOffset(r *http.Request) (int, int, error){
	q := r.URL.Query()

	limit := 20

	if v:=q.Get("limit"); v != ""{
		n, err := strconv.Atoi(v)
		if err != nil || n <= 0 {
			return 0, 0, errors.New("invalid limit")
		}
		if n > 100{
			n = 100
		}
		limit = n
	}
	offset := 0
	if v := q.Get("offset"); v != ""{
		n, err := strconv.Atoi(v)
		if err != nil || n< 0{
			return 0, 0, errors.New("invalid offset")
		}
		offset = n
	}

	return limit, offset, nil
}