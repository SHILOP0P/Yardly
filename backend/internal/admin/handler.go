package admin

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/SHILOP0P/Yardly/backend/internal/auth"
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

	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"users":  users,
		"limit":  limit,
		"offset": offset,
		"q":      q,
	})
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err!=nil || userID <= 0 {
		httpx.WriteError(w, http.StatusBadRequest, "invalid user id")
		return
	}
	u, err := h.repo.GetUser(r.Context(), userID)
	if err!=nil{
		httpx.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, u)
}


func (h *Handler) PatchUser(w http.ResponseWriter, r *http.Request){
	actorAdminID, ok := auth.UserIDFromContext(r.Context())
	if !ok{
		httpx.WriteError(w, http.StatusUnauthorized, "invalid token")
		return
	}

	userID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err!=nil || userID<=0{
		httpx.WriteError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	var req PatchUserRequest
	if err:= json.NewDecoder(r.Body).Decode(&req); err != nil{
		httpx.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if req.Role!=nil{
		if *req.Role==""{
			httpx.WriteError(w, http.StatusBadRequest, "role cannot be empty")
			return
		}
	}
	if req.Ban!=nil && *req.Ban{
		if req.BanReason != nil && *req.BanReason == "" {
			httpx.WriteError(w, http.StatusBadRequest, "ban_reason cannot be empty")
			return
		}
	}

	u, err := h.repo.PatchUser(r.Context(), actorAdminID, userID, req)
	if err!=nil{
		httpx.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, u)
}












//		EVENTS

func (h *Handler) ListAdminEvents(w http.ResponseWriter, r *http.Request){
	var f AdminEventsFilter
	f.Limit = httpx.QueryInt(r, "limit", 20, 1, 100)
	f.Offset = httpx.QueryInt(r, "offset", 0, 0, 1_000_000)

	if s:= r.URL.Query().Get("entity_type"); s!=""{
		f.EntityType = &s
	}
	if s := r.URL.Query().Get("entity_id"); s != "" {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil || v <= 0 {
			httpx.WriteError(w, http.StatusBadRequest, "invalid entity_id")
			return
		}
		f.EntityID = &v
	}

	if s := r.URL.Query().Get("actor_user_id"); s != "" {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil || v <= 0 {
			httpx.WriteError(w, http.StatusBadRequest, "invalid actor_user_id")
			return
		}
		f.ActorUserID = &v
	}

	evs, err := h.repo.ListAdminEvents(r.Context(), f)
	if err!=nil{
		httpx.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"events": evs,
		"limit":  f.Limit,
		"offset": f.Offset,
	})
}


func (h *Handler) ListBookings(w http.ResponseWriter, r *http.Request){
	q := r.URL.Query()
	
	var f AdminBookingsFilter
	limit := httpx.QueryInt(r, "limit", 20, 1, 100)
	offset := httpx.QueryInt(r, "offset", 0, 0, 1_000_000)

	f.Limit = limit
	f.Offset = offset
	if s := q.Get("status"); s != "" {
		f.Status = &s
	}
	if s := q.Get("type"); s != "" {
		f.Type = &s
	}
	if s := q.Get("item_id"); s != "" {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil || v <= 0 {
			httpx.WriteError(w, http.StatusBadRequest, "invalid item_id")
			return
		}
		f.ItemID = &v
	}
	if s := q.Get("user_id"); s != "" {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil || v <= 0 {
			httpx.WriteError(w, http.StatusBadRequest, "invalid user_id")
			return
		}
		f.UserID = &v
	}

	list, err:= h.repo.ListBookings(r.Context(), f)
	if err!=nil{
		httpx.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"bookings": list,
		"limit":   limit,
		"offset":  offset,
	})
}


func (h *Handler) GetBooking(w http.ResponseWriter, r *http.Request) {
	bookingID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || bookingID <= 0 {
		httpx.WriteError(w, http.StatusBadRequest, "invalid booking id")
		return
	}

	b, err := h.repo.GetBooking(r.Context(), bookingID)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, b)
}





func (h *Handler) ListBookingEvents(w http.ResponseWriter, r *http.Request){
	bookingID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err!=nil || bookingID<=0{
		httpx.WriteError(w, http.StatusBadRequest, "invalid booking id")
		return
	}
	
	limit := httpx.QueryInt(r, "limit", 20, 1, 100)
	offset := httpx.QueryInt(r, "offset", 0, 0, 1_000_000)

	events, err := h.repo.ListBookingEvents(r.Context(), bookingID, limit, offset)
	if err!=nil{
		httpx.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"events": events,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *Handler) ListItems(w http.ResponseWriter, r *http.Request){
	qp := r.URL.Query()

	limit := httpx.QueryInt(r, "limit", 20, 1, 100)
	offset := httpx.QueryInt(r, "offset", 0, 0, 1_000_000)

	var f AdminItemsFilter
	f.Limit=limit
	f.Offset=offset

	if s:=qp.Get("q"); s!=""{
		f.Q = &s
	}


	if s := qp.Get("status"); s != "" {
		f.Status = &s
	}
	if s := qp.Get("mode"); s != "" {
		f.Mode = &s
	}
	f.IncludeDeleted = parseBool(qp.Get("include_deleted"))
	f.IncludeArchived = parseBool(qp.Get("include_archived"))
	f.IncludeTransferred = parseBool(qp.Get("include_transferred"))

	items, err := h.repo.ListItems(r.Context(), f)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"items":  items,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *Handler) PatchItem(w http.ResponseWriter, r *http.Request){
	actorID, ok := auth.UserIDFromContext(r.Context())
	if !ok{
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	itemID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err!=nil||itemID<=0{
		httpx.WriteError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	var req PatchItemRequest
	if err:=httpx.ReadJSON(r, &req); err !=nil{
		httpx.WriteError(w, http.StatusBadRequest, "invalid item id")
		return
	}
	if req.Title == nil && req.Mode == nil && req.Status == nil {
		httpx.WriteError(w, http.StatusBadRequest, "empty patch")
		return
	}

	it, err := h.repo.PatchItem(r.Context(), actorID, itemID, req)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, it)
}

func (h *Handler) BlockItem(w http.ResponseWriter, r *http.Request) {
	actorID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	itemID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || itemID <= 0 {
		httpx.WriteError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	var body ModerationRequest
	_ = httpx.ReadJSON(r, &body) // reason опционален, поэтому ошибку можно игнорить

	it, err := h.repo.BlockItem(r.Context(), actorID, itemID, body.Reason)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, it)
}

func (h *Handler) UnblockItem(w http.ResponseWriter, r *http.Request) {
	actorID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	itemID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || itemID <= 0 {
		httpx.WriteError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	var body ModerationRequest
	_ = httpx.ReadJSON(r, &body)

	it, err := h.repo.UnblockItem(r.Context(), actorID, itemID, body.Reason)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, it)
}

func (h *Handler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	actorID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	itemID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || itemID <= 0 {
		httpx.WriteError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	var body ModerationRequest
	_ = httpx.ReadJSON(r, &body)

	it, err := h.repo.DeleteItem(r.Context(), actorID, itemID, body.Reason)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, it)
}




//	HELPERS
func parseBool(s string) bool {
	switch s {
	case "1", "true", "True", "TRUE", "yes", "on":
		return true
	default:
		return false
	}
}