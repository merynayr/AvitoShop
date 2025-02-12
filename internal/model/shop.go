package model

import "time"

// Received модель для хранения факта получения монет
type Received struct {
	ID           int64     `json:"id"`
	FromUsername string    `json:"from_user_id"`
	Amount       int64     `json:"amount"`
	CreatedAt    time.Time `json:"created_at"`
}

// Sent модель для хранения факта передачи монет
type Sent struct {
	ID         int64     `json:"id"`
	ToUsername string    `json:"to_user_id"`
	Amount     int64     `json:"amount"`
	CreatedAt  time.Time `json:"created_at"`
}

// InventoryItem модель для хранения факта покупки мерча
type InventoryItem struct {
	ItemName string `json:"item_name"`
	Quantity int64  `json:"quantity"`
}
