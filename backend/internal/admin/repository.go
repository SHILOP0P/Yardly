package admin

import "context"

type Repo interface {
	ListUsers(ctx context.Context, q string, limit, offset int) ([]UserListItem, error)
}
