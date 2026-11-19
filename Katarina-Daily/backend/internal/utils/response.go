package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response is the single consistent JSON envelope used by every
// endpoint in the API, per the project spec.
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// Success writes a 2xx success envelope.
func Success(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SuccessWithMeta writes a 2xx success envelope including pagination
// (or other) metadata.
func SuccessWithMeta(c *gin.Context, status int, message string, data interface{}, meta interface{}) {
	c.JSON(status, Response{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// Fail writes an error envelope with an explicit HTTP status and
// machine-readable error code.
func Fail(c *gin.Context, status int, message string, errorCode string) {
	c.JSON(status, Response{
		Success: false,
		Message: message,
		Error:   errorCode,
	})
}

// FailFromAppError renders an *AppError using its own HTTP status/code.
func FailFromAppError(c *gin.Context, err error) {
	if ae, ok := AsAppError(err); ok {
		Fail(c, ae.Status, ae.Message, ae.Code)
		return
	}
	Fail(c, http.StatusInternalServerError, "internal server error", "INTERNAL_ERROR")
}
