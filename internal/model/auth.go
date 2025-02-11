package model

import "github.com/dgrijalva/jwt-go"

// UserInfo структура базовой информации о пользователе
type UserInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UserClaims структура claims jwt-токена
type UserClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
}
