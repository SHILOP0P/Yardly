package booking

import "time"

type Event struct {
	ID         int64      `json:"id"`
	BookingID  int64      `json:"booking_id"`
	ActorUserID *int64    `json:"actor_user_id,omitempty"`
	Action     string     `json:"action"`
	FromStatus *Status    `json:"from_status,omitempty"`
	ToStatus   *Status    `json:"to_status,omitempty"`
	Meta       any        `json:"meta,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}
