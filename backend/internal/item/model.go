package item

type Status string

const (
	StatusActive   Status = "active"   // показываем в ленте
	StatusArchived Status = "archived" // ушло из оборота (продано/отдано/скрыто)
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
	return s == StatusActive || s == StatusArchived
}

func (m DealMode) Valid() bool {
	switch m {
	case DealSale, DealRent, DealFree, DealSaleRent:
		return true
	default:
		return false
	}
}
