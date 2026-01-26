package booking

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
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
	Start string `json:"start_at"` // RFC3339, например 2026-01-20T10:00:00Z
	End   string `json:"end_at"`   // RFC3339
}

type approveResponseDTO struct {
	Approved Booking   `json:"approved"`
	Declined []Booking `json:"declined"`
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
	start, err1 := time.Parse(time.RFC3339, dto.Start)
	end, err2 := time.Parse(time.RFC3339, dto.End)
	if err1 != nil || err2 != nil {
		log.Println("start/end must be RFC3339", err1, err2)
		httpx.WriteError(w, http.StatusBadRequest, "start/end must be RFC3339")
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
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
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

	b, err := h.repo.ApproveRent(r.Context(), bookingID, ownerID)
	if err != nil {
		writeBookingError(w, "approve", err)
		return
	}

	httpx.WriteJSON(w,http.StatusOK, b)

}

// POST /api/bookings/{id}/return
func (h *Handler) Return(w http.ResponseWriter, r *http.Request) {
	bookingID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil{
		httpx.WriteError(w, http.StatusBadRequest, "invalid booking id")
		return
	}
	ownerID, ok := auth.UserIDFromContext(r.Context())
	if !ok { 
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	b, err := h.repo.ReturnRent(r.Context(), bookingID, ownerID)
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

	ownerID, ok := auth.UserIDFromContext(r.Context())
	if !ok { 
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	b, err := h.repo.HandoverRent(r.Context(), bookingID, ownerID, time.Now().UTC())
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

	b, err := h.repo.CancelRent(r.Context(), bookingID,requesterID)
	if err!= nil{
		writeBookingError(w, "cancel", err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, b)
}


func (h *Handler) ListMyBookings(w http.ResponseWriter, r *http.Request){
	requesterID, ok := auth.UserIDFromContext(r.Context())
	if !ok{
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	q:= r.URL.Query()

	limit:= 20
	if v := q.Get("limit"); v!=""{
		n, err := strconv.Atoi(v)
		if err!= nil||n<=0{
			httpx.WriteError(w, http.StatusBadRequest, "invalid limit")
			return
		}
		if n>100{
			n =100
		}
		limit = n
	}

	offset:=0
	if v:=q.Get("offset"); v!=""{
		n, err := strconv.Atoi(v)
		if err!= nil||n<0{
			httpx.WriteError(w, http.StatusBadRequest, "invalid offset")
			return
		}
		offset = n
	}
	
	statuses, err := parseStatuses(q["status"])
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	out, err := h.repo.ListMyBookings(r.Context(), requesterID,statuses,limit,offset)
	if err!= nil{
		log.Println("list my bookings error:", err)
		httpx.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, out)
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


func parseStatuses(values []string)([]Status, error){
	if len(values) == 0{
		return nil, nil
	}

	out := make([]Status, 0, len(values))
	for _, v := range values{
		switch Status(v){
			case StatusRequested,
			StatusApproved,
			StatusHandoverPending,
			StatusInUse,
			StatusReturnPending,
			StatusCompleted,
			StatusDeclined,
			StatusCanceled,
			StatusExpired:
			out = append(out, Status(v))
		default:
			return nil, errors.New("invalid status: " + v)
		}
	}
	return out, nil
}