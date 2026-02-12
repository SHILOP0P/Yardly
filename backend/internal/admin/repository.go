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

	//Items
	ListItems(ctx context.Context, f AdminItemsFilter) ([]AdminItem, error)
	GetItem(ctx context.Context, id int64) (AdminItem, error)

	PatchItem(ctx context.Context, actorAdminID, itemID int64, req PatchItemRequest) (AdminItem, error)

	BlockItem(ctx context.Context, actorAdminID, itemID int64, reason *string) (AdminItem, error)
	UnblockItem(ctx context.Context, actorAdminID, itemID int64, reason *string) (AdminItem, error)

	DeleteItem(ctx context.Context, actorAdminID, itemID int64, reason *string) (AdminItem, error)

	

	//events
	ListAdminEvents(ctx context.Context, f AdminEventsFilter) ([]AdminEvent, error)

}


