package admin

import (
	"context"
	"github.com/SHILOP0P/Yardly/backend/internal/booking"
)

type Repo interface {
	//User
	ListUsers(ctx context.Context, q string, limit, offset int) ([]UserListItem, error)
	GetUser(ctx context.Context, id int64)(UserListItem, error)
	PatchUser(ctx context.Context, actorAdminID int64, id int64, req PatchUserRequest) (UserListItem, error)

	//Booking
	ListBookings(ctx context.Context, f AdminBookingsFilter) ([]AdminBooking, error)
	GetBooking(ctx context.Context, id int64) (AdminBooking, error)
	ListBookingEvents(ctx context.Context, bookingID int64, limit, offset int) ([]booking.Event, error)






	ListAdminEventsByEntity(ctx context.Context, entityType string, entityID int64, limit, offset int,) ([]AdminEvent, error)


	//events
	ListAdminEvents(ctx context.Context, f AdminEventsFilter) ([]AdminEvent, error)

}


