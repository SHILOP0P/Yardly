package user

import "time"

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
	RoleSuperAdmin Role = "superadmin"
)

type User struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`

	Role         Role      `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	BannedAt  *time.Time `json:"banned_at,omitempty"`
	BanExpiresAt *time.Time `json:"ban_expires_at,omitempty"`
	BanReason *string    `json:"ban_reason,omitempty"`
}

type Profile struct {
	UserID    int64      `json:"user_id"`
	FirstName string     `json:"first_name"`
	LastName  *string    `json:"last_name,omitempty"`
	BirthDate *time.Time `json:"birth_date,omitempty"`
	Gender    *string    `json:"gender,omitempty"`
	AvatarURL *string    `json:"avatar_url,omitempty"`
	UpdatedAt time.Time  `json:"updated_at"`
}
