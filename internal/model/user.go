package model

// User модель в БД
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Coins    int64  `json:"coins"`
}
