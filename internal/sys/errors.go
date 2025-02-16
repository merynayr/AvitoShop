package sys

import "github.com/merynayr/AvitoShop/internal/sys/codes"

// Константы с текстами ошибок
const (
	ErrItemNotFound            = "item not found"
	ErrNotEnoughCoins          = "not enough coins"
	ErrSelfTransferNotAllowed  = "you can't transfer money to yourself"
	ErrInvalidRefreshToken     = "invalid refresh token"
	ErrInvalidPassword         = "invalid password"
	ErrAuthHeaderNotProvided   = "authorization header is not provided"
	ErrInvalidAuthHeaderFormat = "invalid authorization header format"
	ErrInvalidAccessToken      = "access token is invalid"
	ErrAccessDenied            = "access denied"
	ErrInvalidRequest          = "invalid request"
	ErrInvalidUser             = "invalid user"
	ErrUserNotFound            = "user not found"
	ErrRecipientNotFound       = "recipient not found"
)

// Готовые объекты ошибок с комментариями
var (
	// ItemNotFoundError возникает, когда запрашиваемый объект не найден. Код ошибки: 404 (Not Found)
	ItemNotFoundError = NewCommonError(ErrItemNotFound, codes.NotFound)

	// NotEnoughCoinsError возникает, когда у пользователя недостаточно монет для перевода. Код ошибки: 400 (Bad Request)
	NotEnoughCoinsError = NewCommonError(ErrNotEnoughCoins, codes.BadRequest)

	// SelfTransferNotAllowedError возникает, когда пользователь пытается перевести монеты самому себе. Код ошибки: 400 (Bad Request)
	SelfTransferNotAllowedError = NewCommonError(ErrSelfTransferNotAllowed, codes.BadRequest)

	// InvalidRefreshTokenError возникает, когда передан недействительный refresh-токен. Код ошибки: 401 (Unauthorized)
	InvalidRefreshTokenError = NewCommonError(ErrInvalidRefreshToken, codes.Unauthorized)

	// InvalidPasswordError возникает, когда введён неверный пароль. Код ошибки: 401 (Unauthorized)
	InvalidPasswordError = NewCommonError(ErrInvalidPassword, codes.Unauthorized)

	// AuthHeaderNotProvidedError возникает, когда отсутствует заголовок авторизации. Код ошибки: 401 (Unauthorized)
	AuthHeaderNotProvidedError = NewCommonError(ErrAuthHeaderNotProvided, codes.Unauthorized)

	// InvalidAuthHeaderFormatError возникает, когда формат заголовка авторизации неверный. Код ошибки: 401 (Unauthorized)
	InvalidAuthHeaderFormatError = NewCommonError(ErrInvalidAuthHeaderFormat, codes.Unauthorized)

	// InvalidAccessTokenError возникает, когда передан недействительный access-токен. Код ошибки: 401 (Unauthorized)
	InvalidAccessTokenError = NewCommonError(ErrInvalidAccessToken, codes.Unauthorized)

	// AccessDeniedError возникает, когда у пользователя недостаточно прав для выполнения операции. Код ошибки: 403 (Forbidden)
	AccessDeniedError = NewCommonError(ErrAccessDenied, codes.Forbidden)

	// InvalidRequestError возникает, когда передан некорректный запрос. Код ошибки: 400 (Bad Request)
	InvalidRequestError = NewCommonError(ErrInvalidRequest, codes.BadRequest)

	// InvalidUserError возникает, когда передан недействительный пользователь. Код ошибки: 400 (Bad Request)
	InvalidUserError = NewCommonError(ErrInvalidUser, codes.BadRequest)

	// UserNotFoundError возникает, когда запрашиваемый пользователь не найден. Код ошибки: 404 (Not Found)
	UserNotFoundError = NewCommonError(ErrUserNotFound, codes.NotFound)

	// RecipientNotFoundError возникает, когда запрашиваемый получатель перевода не найден. Код ошибки: 404 (Not Found)
	RecipientNotFoundError = NewCommonError(ErrRecipientNotFound, codes.NotFound)
)
