package model

// User модель пользователя в репо слое
type User struct {
	ID       int64  `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
	Coins    int64  `db:"coins"`
}

// UserInfo модель пользователя для Авторизации
type UserInfo struct {
	Username string `db:"username"`
	Password string `db:"password"`
}
