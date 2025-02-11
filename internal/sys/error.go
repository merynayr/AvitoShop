package sys

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/merynayr/AvitoShop/internal/sys/codes"
)

// Структура для обработчки http ошибок
type commonError struct {
	msg  string
	code codes.Code
}

// NewCommonError создаёт новвую ошибку
func NewCommonError(msg string, code codes.Code) *commonError {
	return &commonError{msg, code}
}

func (r *commonError) Error() string {
	return r.msg
}

func (r *commonError) Code() codes.Code {
	return r.code
}

// IsCommonError проверяет на соответствие ошибке
func IsCommonError(err error) bool {
	var ce *commonError
	return errors.As(err, &ce)
}

// GetCommonError получает ошбику
func GetCommonError(err error) *commonError {
	var ce *commonError
	if !errors.As(err, &ce) {
		return nil
	}

	return ce
}

// HandleError обрабатывает ошибки и отправляет корректный HTTP-ответ
func HandleError(c *gin.Context, err error) {
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
