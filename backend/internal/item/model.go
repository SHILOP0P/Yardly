package item

import(
	"time"
)

type Status string

const (
	StatusActive   Status = "active"   // показываем в ленте
	StatusInUse    Status = "in_use"
	StatusArchived Status = "archived"
	StatusDeleted  Status = "deleted"
	StatusTransferred Status = "transferred"
)

type DealMode string

const (
	DealSale     DealMode = "sale"       // продам
	DealRent     DealMode = "rent"       // сдам в аренду
	DealFree     DealMode = "free"       // отдам бесплатно
	DealSaleRent DealMode = "sale_rent"  // продам или сдам
)

type ItemImage struct {
	ID        int64     `json:"id"`
	ItemID    int64     `json:"item_id"`
	URL       string    `json:"url"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
}

type Item struct {
	ID      int64    `json:"id"`
	OwnerID int64    `json:"owner_id"`
	Title   string   `json:"title"`
	Status  Status   `json:"status"`
	Mode    DealMode `json:"mode"`

	Description string `json:"description,omitempty"`
	Price       int64  `json:"price,omitempty"`   // копейки/центы, чтобы без float
	Deposit     int64  `json:"deposit,omitempty"`
	Location    string `json:"location,omitempty"`
	Category    string `json:"category,omitempty"`

	Images []ItemImage `json:"images,omitempty"`
}



func (s Status) Valid() bool {
	switch s{
	case StatusActive, StatusArchived, StatusDeleted, StatusInUse, StatusTransferred:
		return true
	default:
		return false
	}
	
}

func (m DealMode) Valid() bool {
	switch m {
	case DealSale, DealRent, DealFree, DealSaleRent:
		return true
	default:
		return false
	}
}
