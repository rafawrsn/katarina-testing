package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"fashion-ecommerce-backend/internal/utils"
)

// Recovery converts panics into a standard JSON error response instead
// of crashing the process or leaking a stack trace to the client.
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		slog.Error("panic recovered", "error", recovered, "path", c.Request.URL.Path)
		utils.Fail(c, http.StatusInternalServerError, "internal server error", "INTERNAL_ERROR")
		c.Abort()
	})
}
