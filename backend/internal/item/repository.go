package item

import "context"

type ListFilter struct {
	Status []Status
	Mode   *DealMode
	Limit  int
	Offset int
}

type Repo interface {
	Create(ctx context.Context, it *Item) error
	List(ctx context.Context, f ListFilter) ([]Item, error)
	GetByID(ctx context.Context, id int64) (Item, error)

	ListByOwnerPublic(ctx context.Context, ownerID int64, limit, offset int)([]Item, error)
	ListMyItems(ctx context.Context, ownerId int64, limit, offset int)([]Item, error)

	// Images
	ListImages(ctx context.Context, itemID int64) ([]ItemImage, error)
	AddImage(ctx context.Context, itemID int64, url string) (ItemImage, error)
	DeleteImage(ctx context.Context, itemID int64, imageID int64) error
}
 