package booking

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"context"

	"github.com/SHILOP0P/Yardly/backend/internal/auth"
	"github.com/SHILOP0P/Yardly/backend/internal/httpx"
	"github.com/SHILOP0P/Yardly/backend/internal/item"
)

type Handler struct{
	repo Repo
	items ItemGetter
}

type ItemGetter interface{
	GetByID(ctx context.Context, id int64)(item.Item, error)
}

func NewHandler(repo Repo, items ItemGetter) *Handler{
	return &Handler{repo: repo, items: items}
}

type createRentRequestDTO struct{
	Start string `json:"start_at"` // "2006-01-02", например 2026-01-20T10:00:00Z
	End   string `json:"end_at"`   // "2006-01-02"
}

type approveResponseDTO struct {
	Approved Booking   `json:"approved"`
	Declined []Booking `json:"declined"`
}

type DayRange struct {
	Start string `json:"start"` // YYYY-MM-DD
	End   string `json:"end"`   // YYYY-MM-DD
}

type availabilityResponse struct {
	ItemID      int64            `json:"item_id"`
	From        string           `json:"from"`
	To          string           `json:"to"`
	Timezone    string           `json:"timezone"`
	IsInUseNow  bool             `json:"is_in_use_now"`
	Busy        []DayRange       `json:"busy"`
}


type createBookingRequestDTO struct {
	Type  string `json:"type"` // "rent" | "buy" | "give"
	Start string `json:"start_at,omitempty"`
	End   string `json:"end_at,omitempty"`
}

type upcomingByItemResponse struct {
	ItemID         int64    `json:"item_id"`
	IsInUse        bool     `json:"is_in_use"`
	CurrentBooking *Booking `json:"current_booking,omitempty"`
	Upcoming       []Booking `json:"upcoming"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request){
	itemID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err!=nil{
		log.Println("Invalid item ID:", err)
		httpx.WriteError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	var dto createBookingRequestDTO
	if err:=json.NewDecoder(r.Body).Decode(&dto); err!= nil{
		httpx.WriteError(w, http.StatusBadRequest, "invalid json body")
		return
	}
	
	tp:=Type(strings.TrimSpace(dto.Type))
	if tp==""{
		httpx.WriteError(w, http.StatusBadRequest, "type is required")
		return
	}
	switch tp {
	case TypeRent, TypeBuy, TypeGive:
	default:
		httpx.WriteError(w, http.StatusBadRequest, "invalid type")
		return
	}

	requesterID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	it, err := h.items.GetByID(r.Context(), itemID)
	if err != nil {
		httpx.WriteError(w, http.StatusNotFound, "item not found")
		return
	}

	switch it.Status {
	case "active", "in_use":
		// ok
	case "archived", "deleted", "transferred":
		httpx.WriteError(w, http.StatusConflict, "item is not available")
		return
	default:
		httpx.WriteError(w, http.StatusConflict, "item is not available")
		return
	}


	ownerID := it.OwnerID
	if ownerID == requesterID {
		httpx.WriteError(w, http.StatusBadRequest, "cannot book your own item")
		return
	}

	b := Booking{
		ItemID:      itemID,
		RequesterID: requesterID,
		OwnerID:     ownerID,
		Type:        tp,
		Status:      StatusRequested,
	}

	if tp == TypeRent{
		startDay, err1 := time.Parse("2006-01-02", dto.Start)
		endDay, err2 := time.Parse("2006-01-02", dto.End)
		if err1 != nil || err2 != nil {
			httpx.WriteError(w, http.StatusBadRequest, "start/end must be YYYY-MM-DD")
			return
		}
		if endDay.Before(startDay) {
			httpx.WriteError(w, http.StatusBadRequest, "end must be >= start")
			return
		}
		start := time.Date(startDay.Year(), startDay.Month(), startDay.Day(), 0, 0, 0, 0, time.UTC)
		endExclusive := time.Date(endDay.Year(), endDay.Month(), endDay.Day(), 0, 0, 0, 0, time.UTC).Add(24 * time.Hour)

		b.Start = &start
		b.End = &endExclusive
	} else {
		if strings.TrimSpace(dto.Start)!=""||strings.TrimSpace(dto.End)!=""{
			httpx.WriteError(w, http.StatusBadRequest, "start/end are not allowed for buy/give")
			return
		}
	}
	if err := h.repo.Create(r.Context(), &b); err != nil {
		switch {
		case errors.Is(err, ErrDuplicateActiveRequest):
			httpx.WriteError(w, http.StatusConflict, "active request already exists")
		default:
			log.Println("create booking error:", err)
			httpx.WriteError(w, http.StatusInternalServerError, "internal error")
		}
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, b)
}




func (h *Handler) CreateRent(w http.ResponseWriter, r *http.Request){
	itemID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err!=nil{
		log.Println("Invalid item ID:", err)
		httpx.WriteError(w, http.StatusBadRequest, "invalid item id")
		return
	}
	
	var dto createRentRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err!= nil {
		log.Println("invalid json body:", err)
		httpx.WriteError(w, http.StatusBadRequest, "invalid json body")
		return
	}
	start, err1 := time.Parse("2006-01-02", dto.Start)
	end, err2 := time.Parse("2006-01-02", dto.End)
	if err1 != nil || err2 != nil {
		log.Println("start/end must be '2006-01-02'", err1, err2)
		httpx.WriteError(w, http.StatusBadRequest, "start/end must be '2006-01-02'")
		return
	}
	if !start.Before(end) {
		log.Println("start must be before end")
		httpx.WriteError(w, http.StatusBadRequest, "start must be before end")
		return
	}

	requesterID, ok := auth.UserIDFromContext(r.Context())
	if !ok{
		httpx.WriteError(w, http.StatusUnauthorized, "requesterID not found unauthorized")
		return
	}
	it, err := h.items.GetByID(r.Context(), itemID)
	if err != nil{
		httpx.WriteError(w, http.StatusNotFound, "item not found")
		return
	}
	ownerID := it.OwnerID
	if ownerID == requesterID{
		httpx.WriteError(w, http.StatusBadRequest, "cannot book your own item")
		return
	}


	b := Booking{
		ItemID:      itemID,
		RequesterID: requesterID,
		OwnerID:     ownerID,
		Type:        TypeRent,
		Status:      StatusRequested,
		Start:       &start,
		End:         &end,
	}


	if err:=h.repo.Create(r.Context(), &b); err != nil{
		switch {
		case errors.Is(err, ErrDuplicateActiveRequest):
			httpx.WriteError(w, http.StatusConflict, "active request already exists")
		default:
			log.Println("create rent error:", err)
			httpx.WriteError(w, http.StatusInternalServerError, "internal error")
		}
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, b)
}

func (h *Handler) ListBusyForItem(w http.ResponseWriter, r *http.Request){
	itemID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err!=nil{
		httpx.WriteError(w, http.StatusBadRequest,"invalid item id")
		return
	}

	bookings, err := h.repo.ListByItem(r.Context(),itemID)
	if err!=nil{
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	busy := make([]Booking,0,len(bookings))
	for _, b:= range bookings{
		if b.Type!=TypeRent{
			continue
		}
		if b.Start==nil||b.End==nil{
			continue
		}
		switch b.Status{
		case StatusApproved, StatusHandoverPending, StatusInUse,StatusReturnPending:
			busy = append(busy, b)
		}
	}
	httpx.WriteJSON(w, http.StatusOK, busy)
}


// POST /api/bookings/{id}/approve
func (h *Handler) Approve(w http.ResponseWriter, r *http.Request) {
	bookingID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err!=nil{
		httpx.WriteError(w, http.StatusBadRequest, "invalid booking id")
		return
	}
	ownerID, ok := auth.UserIDFromContext(r.Context())
	if !ok { 
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	b0, err := h.repo.GetByID(r.Context(), bookingID)
	if err != nil {
		writeBookingError(w, "get booking", err)
		return
	}
	if b0.OwnerID != ownerID {
		httpx.WriteError(w, http.StatusForbidden, "forbidden")
		return
	}

	var b Booking
	switch b0.Type{
	case TypeRent:
		b, err = h.repo.ApproveRent(r.Context(),bookingID, ownerID)
	case TypeBuy, TypeGive:
		b, err = h.repo.ApproveTransfer(r.Context(), bookingID, ownerID, time.Now().UTC())
	default:
		httpx.WriteError(w, http.StatusBadRequest, "invalid type")
		return
	}
	
	if err != nil {
		writeBookingError(w, "approve", err)
		return
	}


	httpx.WriteJSON(w,http.StatusOK, b)

}


func (h *Handler) Return(w http.ResponseWriter, r *http.Request) {
	bookingID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil{
		httpx.WriteError(w, http.StatusBadRequest, "invalid booking id")
		return
	}
	actorID, ok := auth.UserIDFromContext(r.Context())
	if !ok { 
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	b, err := h.repo.ReturnRent(r.Context(), bookingID, actorID, time.Now().UTC())
	if err != nil {
		writeBookingError(w, "return", err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, b)
}

func (h *Handler) Handover(w http.ResponseWriter, r *http.Request) {
	bookingID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid booking id")
		return
	}

	actorID, ok := auth.UserIDFromContext(r.Context())
	if !ok { 
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	b0, err := h.repo.GetByID(r.Context(), bookingID)
	if err != nil {
		writeBookingError(w, "get booking", err)
		return
	}

	var b Booking
	switch b0.Type {
	case TypeRent:
		b, err = h.repo.HandoverRent(r.Context(), bookingID, actorID, time.Now().UTC())
	case TypeBuy, TypeGive:
		b, err = h.repo.HandoverTransfer(r.Context(), bookingID, actorID, time.Now().UTC())
	default:
		httpx.WriteError(w, http.StatusBadRequest, "invalid type")
		return
	}
	if err != nil {
		writeBookingError(w, "handover", err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, b)
}


func (h *Handler) Cancel(w http.ResponseWriter, r *http.Request){
	bookingID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || bookingID <= 0 {
		httpx.WriteError(w, http.StatusBadRequest, "invalid booking id")
		return
	}
	requesterID, ok := auth.UserIDFromContext(r.Context())
	if !ok{
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	b0, err := h.repo.GetByID(r.Context(), bookingID)
	if err != nil {
		writeBookingError(w, "get booking", err)
		return
	}
	if b0.RequesterID != requesterID {
		httpx.WriteError(w, http.StatusForbidden, "forbidden")
		return
	}

	var b Booking
	switch b0.Type {
	case TypeRent:
		b, err = h.repo.CancelRent(r.Context(), bookingID, requesterID)
	case TypeBuy, TypeGive:
		b, err = h.repo.CancelTransfer(r.Context(), bookingID, requesterID)
	default:
		httpx.WriteError(w, http.StatusBadRequest, "invalid type")
		return
	}
	if err != nil {
		writeBookingError(w, "cancel", err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, b)
}


func (h *Handler) ListMyBookings(w http.ResponseWriter, r *http.Request){
	requesterID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	statuses, err := parseStatuses(r)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid status")
		return
	}

	limit, offset, err := parseLimitOffset(r)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid limit/offset")
		return
	}

	out, err := h.repo.ListMyBookings(r.Context(), requesterID,statuses,limit,offset)
	if err!= nil{
		log.Println("list my bookings error:", err)
		httpx.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"items":  out,
		"limit":  limit,
		"offset": offset,
	})
}


func (h *Handler) ListMyItemsBookings(w http.ResponseWriter, r *http.Request){
	ownerID, ok := auth.UserIDFromContext(r.Context())
	if !ok{
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	statuses, err := parseStatuses(r)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid status")
		return
	}

	limit, offset, err := parseLimitOffset(r)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid limit/offset")
		return
	}

	out, err := h.repo.ListMyItemsBookings(r.Context(), ownerID,statuses,limit,offset)
	if err!= nil{
		log.Println("list my items bookings error:", err)
		httpx.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"items":  out,
		"limit":  limit,
		"offset": offset,
	})
}


func (h *Handler) ListEvents(w http.ResponseWriter, r *http.Request){
	bookingID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || bookingID <= 0 {
		httpx.WriteError(w, http.StatusBadRequest, "invalid booking id")
		return
	}
	actorID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	

	b, err := h.repo.GetByID(r.Context(), bookingID)
	if err!= nil {
		writeBookingError(w, "get booking", err)
		return
	}
	if actorID != b.OwnerID && actorID!= b.RequesterID{
		httpx.WriteError(w, http.StatusForbidden, "forbidden")
		return
	}
	
	q:=r.URL.Query()

	limit := 50
	if v:=q.Get("limit");v!=""{
		n, err := strconv.Atoi(v)
		if err != nil{
			httpx.WriteError(w, http.StatusBadRequest, "invalid limit")
			return
		}
		if n > 200{
			n = 200
		}
		limit = n
	}

	offset := 0
	if v:=q.Get("offset"); v!=""{
		n, err := strconv.Atoi(v)
		if err!=nil{
			httpx.WriteError(w, http.StatusBadRequest, "invalid offset")
			return
		}
		offset = n
	}

	events, err := h.repo.ListEvents(r.Context(), bookingID, limit, offset)
	if err != nil{
		log.Println("list events error:", err)
		httpx.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, events)
}

func (h *Handler) UpcomingByItem(w http.ResponseWriter, r *http.Request){
	itemID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || itemID <= 0 {
		httpx.WriteError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	now:=time.Now().UTC()
	cur, upcoming, err := h.repo.ListUpcomingByItem(r.Context(), itemID, now, 20)
	if err!= nil {
		log.Println("upcomingByItem error:", err)
		httpx.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}
	resp := upcomingByItemResponse{
		ItemID: itemID,
		IsInUse: cur!=nil,
		Upcoming: upcoming,
	}
	if cur != nil{
		resp.CurrentBooking = cur
	}
	httpx.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) AvailabilityByItem(w http.ResponseWriter, r *http.Request){
	itemID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || itemID <=0{
		httpx.WriteError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	fromS:= r.URL.Query().Get("from")
	toS := r.URL.Query().Get("to")
	if fromS == "" || toS == ""{
		httpx.WriteError(w, http.StatusBadRequest, "from and to are required (YYYY-MM-DD)")
		return
	}

	fromDay, err1 := time.Parse("2006-01-02", fromS)
	toDay, err2 := time.Parse("2006-01-02", toS)
	if err1 != nil || err2 != nil || !fromDay.Before(toDay) {
		httpx.WriteError(w, http.StatusBadRequest, "invalid from/to range")
		return
	}

	if toDay.Sub(fromDay) > 180*24*time.Hour {
		httpx.WriteError(w, http.StatusBadRequest, "range too large (max 180 days)")
		return
	}

	busy, inUseNow, err:= h.repo.ListBusyDaysByItem(r.Context(), itemID, fromDay, toDay)
	if err != nil {
		log.Println("availability error:", err)
		httpx.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	resp := availabilityResponse{
		ItemID:     itemID,
		From:       fromS,
		To:         toS,
		Timezone:   "UTC",
		IsInUseNow: inUseNow,
		Busy:       busy,
	}
	httpx.WriteJSON(w, http.StatusOK, resp)
}



func (h *Handler) ListMyItemsBookingRequests(w http.ResponseWriter, r *http.Request){
	ownerID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	types, err := parseTypes(r)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid type")
		return
	}
	if len(types) == 0 {
		types = []Type{TypeRent, TypeBuy, TypeGive}
	}

	limit, offset, err := parseLimitOffset(r)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid limit/offset")
		return
	}

	out, err := h.repo.ListMyItemsBookingRequests(r.Context(), ownerID, types, limit, offset)
	if err != nil {
		log.Println("list my items booking requests error:", err)
		httpx.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"items":  out,
		"limit":  limit,
		"offset": offset,
	})
}



func writeBookingError(w http.ResponseWriter, op string, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		httpx.WriteError(w, http.StatusNotFound, "booking not found")
	case errors.Is(err, ErrForbidden):
		httpx.WriteError(w, http.StatusForbidden, "forbidden")
	case errors.Is(err, ErrInvalidState):
		httpx.WriteError(w, http.StatusConflict, err.Error())
	default:
		log.Println(op, "error:", err)
		httpx.WriteError(w, http.StatusInternalServerError, "internal error")
	}
}


func parseLimitOffset(r *http.Request)(limit, offset int, err error){
	q:= r.URL.Query()
	
	limit = 20
	offset = 0

	if s:=q.Get("limit");s!=""{
		v, e:= strconv.Atoi(s)
		if e!=nil{
			return 0, 0, e
		}
		limit = v
	}

	if s := q.Get("offset"); s != "" {
		v, e := strconv.Atoi(s)
		if e != nil {
			return 0, 0, e
		}
		offset = v
	}
	return limit, offset, nil
}

func parseStatuses(r *http.Request)([]Status, error){
	q:=r.URL.Query()

	raw:=q["status"]
	if len(raw)==0{
		if s := q.Get("status"); s!=""{
			raw = []string{s}
		}
	}
	if len(raw) == 0{
		return nil, nil
	}

	var out []Status
	for _, part := range raw{
		for _, token := range strings.Split(part,","){
			s := strings.TrimSpace(token)
			if s ==""{
				continue
			}
			st := Status(s)
			if !isValidStatus(st){
				return nil, ErrInvalidState
			}
			out = append(out, st)
		}
	}
	return out, nil
}

func parseTypes(r *http.Request) ([]Type, error) {
	q := r.URL.Query()

	raw := q["type"]
	if len(raw) == 0 {
		if s := q.Get("type"); s != "" {
			raw = []string{s}
		}
	}
	if len(raw) == 0 {
		// also allow `types=`
		raw = q["types"]
		if len(raw) == 0 {
			if s := q.Get("types"); s != "" {
				raw = []string{s}
			}
		}
	}
	if len(raw) == 0 {
		return nil, nil
	}

	out := make([]Type, 0, len(raw))
	for _, part := range raw {
		for _, token := range strings.Split(part, ",") {
			s := strings.TrimSpace(token)
			if s == "" {
				continue
			}
			t := Type(s)
			if !isValidType(t) {
				return nil, ErrInvalidState
			}
			out = append(out, t)
		}
	}
	return out, nil
}

func isValidType(t Type) bool {
	switch t {
	case TypeGive, TypeRent, TypeBuy:
		return true
	default:
		return false
	}
}




func isValidStatus(s Status) bool {
	switch s {
	case StatusRequested,
		StatusApproved,
		StatusHandoverPending,
		StatusInUse,
		StatusReturnPending,
		StatusCompleted,
		StatusDeclined,
		StatusCanceled,
		StatusExpired:
		return true
	default:
		return false
	}
}