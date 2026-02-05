package favorite

import "time"

type Favorite struct {
	UserID    int64     `json:"user_id"`
	ItemID    int64     `json:"item_id"`
	CreatedAt time.Time `json:"created_at"`
}


type FavoriteItem struct {
	ItemID     int64     `json:"item_id"`
	Title      string    `json:"title"`
	Status     string    `json:"status"`
	Mode       string    `json:"mode"`
	OwnerID    int64     `json:"owner_id"`
	FavoritedAt time.Time `json:"favorited_at"`
}