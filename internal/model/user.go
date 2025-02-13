package model

// User модель в БД
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Coins    int64  `json:"coins"`
}

// UserUpdate модель обновления пользователя сервисного слоя
type UserUpdate struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Coins    int64  `json:"coins"`
}

// CoinHistory модель истории транзакций
type CoinHistory struct {
	Received []Received `json:"received"`
	Sent     []Sent     `json:"sent"`
}

// InfoResponse модель для ответа запроса
type InfoResponse struct {
	Coins       int64           `json:"coins"`
	Inventory   []InventoryItem `json:"inventory"`
	CoinHistory CoinHistory     `json:"coinHistory"`
}

// SendCoinRequest модель запроса
type SendCoinRequest struct {
	ToUser string `json:"toUser" binding:"required"`
	Amount int64  `json:"amount" binding:"required,gt=0"`
}
