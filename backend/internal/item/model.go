package item

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

type Item struct {
	ID     int64  `json:"id"`
	OwnerID int64 `json:"owner_id"`
	Title  string `json:"title"`
	Status Status `json:"status"`
	Mode   DealMode `json:"mode"`
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
