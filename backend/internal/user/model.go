package user

import "time"

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type User struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         Role      `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
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
