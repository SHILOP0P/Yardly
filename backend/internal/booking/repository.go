package booking

import ("context"
		"time"
	)

type Repo interface{
	Create(ctx context.Context, b *Booking) error

	GetByID(ctx context.Context, id int64) (Booking, error)

	ListByItem(ctx context.Context, itemID int64) ([]Booking, error)
	ListMyBookings(ctx context.Context, requesterID int64, statuses[]Status, limit, offset int)([]Booking, error)
	ListMyItemsBookings(ctx context.Context, ownerID int64, statuses[]Status, limit, offset int)([]Booking, error)
	ListMyItemsBookingRequests(ctx context.Context, ownerID int64, types []Type, limit, offset int) ([]Booking, error)


	ListUpcomingByItem(ctx context.Context, itemID int64, now time.Time, limit int) (inUse *Booking, upcoming []Booking, err error)
	ListBusyDaysByItem(ctx context.Context, itemID int64, fromDay, toDay time.Time) ([]DayRange, bool, error)

	ApproveRent(ctx context.Context, bookingID int64, ownerID int64)(Booking, error)
	ReturnRent(ctx context.Context, bookingID int64, actorID int64, now time.Time)(Booking, error)
	HandoverRent(ctx context.Context, bookingID int64, actorID int64, now time.Time) (Booking, error)
	
	ExpireOverdueHandovers(ctx context.Context, now time.Time) (int64, error)
	CancelRent(ctx context.Context, bookingID, requesterID int64) (Booking, error)

	ListEvents(ctx context.Context, bookingID int64, limit, offset int) ([]Event, error)

	//Transfer
	ApproveTransfer(ctx context.Context, bookingID int64, ownerID int64, now time.Time) (Booking, error)
	HandoverTransfer(ctx context.Context, bookingID int64, actorID int64, now time.Time) (Booking, error)
	CancelTransfer(ctx context.Context, bookingID, requesterID int64) (Booking, error)


}