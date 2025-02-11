package model

import "time"

// Transaction модель для хранения факта передачи монет
type Transaction struct {
	ID         int64     `json:"id"`
	FromUserID int64     `json:"from_user_id"`
	ToUserID   int64     `json:"to_user_id"`
	Amount     int64     `json:"amount"`
	CreatedAt  time.Time `json:"created_at"`
}

// InventoryItem модель для хранения факта покупки мерча
type InventoryItem struct {
	ID       int64  `json:"id"`
	UserID   int64  `json:"user_id"`
	ItemName string `json:"item_name"`
	Quantity int64  `json:"quantity"`
}
