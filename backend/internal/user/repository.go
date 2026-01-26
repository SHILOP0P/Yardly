package user

import "context"

type Repo interface {
	CreateWithProfile(ctx context.Context, u *User, p *Profile) error
	GetByEmail(ctx context.Context, email string) (User, error)
	GetByID(ctx context.Context, id int64) (User, Profile, error)
}
