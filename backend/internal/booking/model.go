package booking

import (
	"time"
)

type Type string

const (
	TypeGive Type = "give"
	TypeRent Type = "rent"
	TypeBuy  Type = "buy"
)

type Status string

const (
	StatusRequested      Status = "requested"       // запрос отправлен владельцу
	StatusApproved       Status = "approved"        // владелец подтвердил (интервал занят)
	StatusHandoverPending Status = "handover_pending" // пришла дата start, ждём подтверждения передачи
	StatusInUse          Status = "in_use"          // передача подтверждена, вещь у пользователя
	StatusReturnPending  Status = "return_pending"  // пришла дата end, ждём подтверждения возврата
	StatusCompleted      Status = "completed"       // всё завершено

	StatusDeclined  Status = "declined"
	StatusCanceled Status = "cancelled"
	StatusExpired   Status = "expired" // не подтвердили вовремя (например передачу) — интервал освободился
)


type Booking struct{
	ID          int64     `json:"id"`
	ItemID      int64     `json:"item_id"`
	RequesterID int64     `json:"requester_id"`
	OwnerID     int64     `json:"owner_id"`

	Type   Type   `json:"type"`
	Status Status `json:"status"`

	Start *time.Time `json:"start,omitempty"`
	End   *time.Time `json:"end,omitempty"`

	HandoverDeadline *time.Time `json:"handover_deadline,omitempty"`

	HandoverConfirmedByOwnerAt     *time.Time `json:"handover_confirmed_by_owner_at,omitempty"`
	HandoverConfirmedByRequesterAt *time.Time `json:"handover_confirmed_by_requester_at,omitempty"`
	ReturnConfirmedByOwnerAt       *time.Time `json:"return_confirmed_by_owner_at,omitempty"`
	ReturnConfirmedByRequesterAt   *time.Time `json:"return_confirmed_by_requester_at,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}
