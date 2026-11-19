package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"fashion-ecommerce-backend/internal/models"
	"fashion-ecommerce-backend/internal/utils"
)

// RequireAdmin must run after Auth. It rejects any request whose JWT
// role claim is not "admin".
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if GetUserRole(c) != models.RoleAdmin {
			utils.Fail(c, http.StatusForbidden, "admin access required", "FORBIDDEN")
			c.Abort()
			return
		}
		c.Next()
	}
}
