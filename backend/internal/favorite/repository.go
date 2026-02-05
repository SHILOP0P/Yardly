package favorite

import "context"

type Repo interface {
	Add(ctx context.Context, userID, itemID int64) (Favorite, error)
	Remove(ctx context.Context, userID, itemID int64) error
	List(ctx context.Context, userID int64, limit, offset int) ([]FavoriteItem, error)
	IsFavorite(ctx context.Context, userID, itemID int64) (bool, error)
}
