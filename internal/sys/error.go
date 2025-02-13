package sys

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/merynayr/AvitoShop/internal/logger"
	"github.com/merynayr/AvitoShop/internal/sys/codes"
)

// ErrorResponse Структура для обработчки http ошибок
type ErrorResponse struct {
	msg  string
	code codes.Code
}

// NewCommonError создаёт новвую ошибку
func NewCommonError(msg string, code codes.Code) *ErrorResponse {
	return &ErrorResponse{msg, code}
}

// Error возврашает сообщение ошибки
func (r *ErrorResponse) Error() string {
	return r.msg
}

// Code возвращает код ошибки
func (r *ErrorResponse) Code() codes.Code {
	return r.code
}

// IsCommonError проверяет на соответствие ошибке
func IsCommonError(err error) bool {
	var ce *ErrorResponse
	return errors.As(err, &ce)
}

// GetCommonError получает ошбику
func GetCommonError(err error) *ErrorResponse {
	var ce *ErrorResponse
	if !errors.As(err, &ce) {
		return nil
	}

	return ce
}

// HandleError обрабатывает ошибки и отправляет корректный HTTP-ответ
func HandleError(c *gin.Context, err error) {
	logger.Debug(err.Error())
	if ce := GetCommonError(err); ce != nil {
		c.JSON(int(ce.Code()), gin.H{
			"error": ce.Error(),
			"code":  ce.Code(),
		})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{
		"error": "internal server error",
		"code":  http.StatusInternalServerError,
	})
}
