package admin

import "time"

type UserListItem struct {
	ID           int64      `json:"id"`
	Email        string     `json:"email"`
	Role         string     `json:"role"`
	BannedAt     *time.Time `json:"banned_at,omitempty"`
	BanExpiresAt *time.Time `json:"ban_expires_at,omitempty"`
	BanReason    *string    `json:"ban_reason,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type PatchUserRequest struct {
	Role         *string    `json:"role,omitempty"`
	Ban          *bool      `json:"ban,omitempty"`
	BanReason    *string    `json:"ban_reason,omitempty"`
	BanExpiresAt *time.Time `json:"ban_expires_at,omitempty"`
}


type AdminEvent struct {
	ID        int64     `json:"id"`
	ActorID   int64     `json:"actor_user_id"`
	EntityType string   `json:"entity_type"`
	EntityID  int64     `json:"entity_id"`
	Action    string    `json:"action"`
	Reason    *string   `json:"reason,omitempty"`
	Meta      any       `json:"meta,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type AdminEventsFilter struct {
	EntityType  *string
	EntityID    *int64
	ActorUserID *int64
	Limit       int
	Offset      int
}


//DTO
type AdminBooking struct {
	ID          int64     `json:"id"`
	ItemID      int64     `json:"item_id"`
	RequesterID int64     `json:"requester_id"`
	OwnerID     int64     `json:"owner_id"`

	Type   string   `json:"type"`
	Status string `json:"status"`

	Start *time.Time `json:"start,omitempty"`
	End   *time.Time `json:"end,omitempty"`

	HandoverDeadline *time.Time `json:"handover_deadline,omitempty"`

	HandoverConfirmedByOwnerAt     *time.Time `json:"handover_confirmed_by_owner_at,omitempty"`
	HandoverConfirmedByRequesterAt *time.Time `json:"handover_confirmed_by_requester_at,omitempty"`
	ReturnConfirmedByOwnerAt       *time.Time `json:"return_confirmed_by_owner_at,omitempty"`
	ReturnConfirmedByRequesterAt   *time.Time `json:"return_confirmed_by_requester_at,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}

type AdminBookingsFilter struct {
	Status *string
	Type   *string
	ItemID *int64
	UserID *int64 // requester OR owner

	Limit  int
	Offset int
}

type AdminBookingEvent struct {
	ID          int64      `json:"id"`
	BookingID   int64      `json:"booking_id"`
	ActorUserID *int64     `json:"actor_user_id,omitempty"`
	Action      string     `json:"action"`
	FromStatus  *string    `json:"from_status,omitempty"`
	ToStatus    *string    `json:"to_status,omitempty"`
	Meta        any        `json:"meta,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

type AdminItem struct {
	ID      int64  `json:"id"`
	OwnerID int64  `json:"owner_id"`
	Title   string `json:"title"`
	Status  string `json:"status"`
	Mode    string `json:"mode"`

	BlockedAt   *time.Time `json:"blocked_at,omitempty"`
	BlockReason *string    `json:"block_reason,omitempty"`
}

type AdminItemsFilter struct {
	Q      *string // поиск по title
	Status *string // точный статус
	Mode   *string // точный mode

	IncludeDeleted     bool
	IncludeTransferred bool
	IncludeArchived    bool

	Limit  int
	Offset int
}

type PatchItemRequest struct {
	Title  *string `json:"title,omitempty"`
	Mode   *string `json:"mode,omitempty"`
	Status *string `json:"status,omitempty"` // разрешим, но валидируем
}

type ModerationRequest struct {
	Reason *string `json:"reason,omitempty"`
}
