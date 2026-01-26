package item

import (
	"context"
	"sync"
)

type MemoryRepo struct {
	mu sync.Mutex

	items map[int64]Item
}

func NewMemoryRepo(seed []Item) *MemoryRepo{
	m := &MemoryRepo{
		items: make(map[int64]Item, len(seed)),
	}

	for _, it := range seed {
		m.items[it.ID] = it
	}
	return m
}

func (m *MemoryRepo) List(ctx context.Context, f ListFilter) ([]Item, error){
	_ = ctx

	limit := f.Limit
	if limit<=0 || limit >= 100{
		limit = 20
	}

	offset := f.Offset
	if offset < 0 {
		offset = 0
	}


	m.mu.Lock()
	defer m.mu.Unlock()

	tmp := make([]Item, 0, len(m.items))
	for _, it := range m.items{
		if f.Status!=nil && it.Status != *f.Status{
			continue
		}
		if f.Mode != nil && it.Mode != *f.Mode {
			continue
		}
		tmp = append(tmp, it)
	}
	
	if offset >=len(tmp){
		return []Item{}, nil
	}
	end:=offset+limit
	if end>len(tmp){
		end = len(tmp)
	}
	return tmp[offset:end], nil
}

func (m *MemoryRepo) GetByID(ctx context.Context, id int64) (Item, error){
	_ = ctx

	m.mu.Lock()
	defer m.mu.Unlock()

	it, ok := m.items[id]
	if !ok{
		return Item{}, ErrNotFound
	}
	return it, nil
}