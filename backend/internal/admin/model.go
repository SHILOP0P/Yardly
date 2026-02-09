package admin

import "time"

type UserListItem struct {
	ID        int64      `json:"id"`
	Email     string     `json:"email"`
	Role      string     `json:"role"`
	BannedAt  *time.Time `json:"banned_at,omitempty"`
	BanReason *string    `json:"ban_reason,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

