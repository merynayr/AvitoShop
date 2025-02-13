package model

import "github.com/dgrijalva/jwt-go"

// AuthRequest структура запроса на аутентификацию
type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse структура ответа с токенами
type AuthResponse struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

// UserClaims структура claims jwt-токена
type UserClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
}
