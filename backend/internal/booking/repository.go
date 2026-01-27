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

	ApproveRent(ctx context.Context, bookingID int64, ownerID int64)(Booking, error)
	ReturnRent(ctx context.Context, bookingID int64, ownerID int64)(Booking, error)
	HandoverRent(ctx context.Context, bookingID int64, ownerID int64, now time.Time) (Booking, error)
	ExpireOverdueHandovers(ctx context.Context, now time.Time) (int64, error)
	CancelRent(ctx context.Context, bookingID, requesterID int64) (Booking, error)

}