package item

import "context"

type ListFilter struct {
	Status *Status
	Mode   *DealMode
	Limit  int
	Offset int
}

type Repo interface {
	Create(ctx context.Context, it *Item) error
	List(ctx context.Context, f ListFilter) ([]Item, error)
	GetByID(ctx context.Context, id int64) (Item, error)
}
 