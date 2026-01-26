package booking

import (
	"errors"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)


type Store struct {
	mu sync.Mutex
	nextID int64
	bookings map[int64]Booking
}


func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{
		nextID: 1,
		bookings: make(map[int64]Booking),
	}
}

func (s *Store) CreateRentRequest(itemID int64, requesterID int64, ownerID int64, start, end time.Time) (Booking, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    // Проверяем на пустые значения
    if start.IsZero() || end.IsZero() {
        return Booking{}, errors.New("start and end must be valid dates")
    }

    id := s.nextID
    s.nextID++

    b := Booking{
        ID:          id,
        ItemID:      itemID,
        OwnerID:     ownerID,
        RequesterID: requesterID,
        Type:        TypeRent,
        Status:      StatusRequested,
        Start:       &start,
        End:         &end,
        CreatedAt:   time.Now(),
    }

    s.bookings[id] = b
    return b, nil
}


func (s *Store) Approve(bookingID int64) (approved Booking, declined []Booking, err error){
	s.mu.Lock()
	defer s.mu.Unlock()

	b, ok := s.bookings[bookingID]
	if !ok {
		return Booking{}, nil, ErrNotFound
	}

	if b.Type != TypeRent || b.Start == nil || b.End ==nil{
		return Booking{}, nil, errors.New("only rent bookings with start/end can be approved")
	}

	b.Status = StatusHandoverPending
	deadline := b.Start.Add(24 * time.Hour)
	b.HandoverDeadline = &deadline


	s.bookings[bookingID] = b
	approved = b
	

	
	for id, other := range s.bookings {
		if id == bookingID{
			continue
		}
		if other.ItemID != b.ItemID{
			continue
		}
		if other.Type!=TypeRent || other.Status != StatusRequested {
			continue
		}

		if overlaps(*b.Start, *b.End, *other.Start, *other.End) {
			other.Status = StatusDeclined
			s.bookings[id] = other
			declined = append(declined, other)
		}
	}
	return approved, declined, nil
}


func (s *Store) ListBusyForItem(itemID int64) []Booking{
	s.mu.Lock()
	defer s.mu.Unlock()

	var res []Booking
	for _, b := range s.bookings{
		if b.ItemID!=itemID{
			continue
		}
		if b.Type != TypeRent{
			continue
		}
		if b.Start ==nil || b.End ==nil{
			continue
		}
		if isExpired(b){
			b.Status = StatusExpired
			s.bookings[b.ID] = b
			continue
		}
		switch b.Status{
		case StatusApproved, StatusHandoverPending, StatusInUse, StatusReturnPending:
			res = append(res, b)
		}
	}
	return res
}

func overlaps (aStart, aEnd, bStart, bEnd time.Time) bool{
	return aStart.Before(bEnd) && aEnd.After(bStart)
}

func isExpired(b Booking) bool {
	if b.HandoverDeadline == nil {
		return false
	}

	if b.Status == StatusInUse || b.Status == StatusCompleted {
		return false
	}
	return time.Now().After(*b.HandoverDeadline)
}


func (s *Store) ConfirmReturn(bookingID int64) (Booking, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	b, ok := s.bookings[bookingID]
	if !ok {
		return Booking{}, ErrNotFound
	}

	// Нельзя возвращать то, что ещё не передали
	if b.Status != StatusInUse {
		return Booking{}, errors.New("return allowed only in in_use")
	}

	b.Status = StatusCompleted
	s.bookings[bookingID] = b
	return b, nil
}
