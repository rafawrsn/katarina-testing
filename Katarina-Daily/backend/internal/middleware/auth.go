package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"fashion-ecommerce-backend/internal/utils"
)

const (
	ctxUserID = "auth_user_id"
	ctxEmail  = "auth_email"
	ctxRole   = "auth_role"
)

// Auth requires a valid "Authorization: Bearer <token>" header and
// stores the decoded user id/email/role on the request context for
// downstream handlers/services to use.
func Auth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			utils.Fail(c, http.StatusUnauthorized, "missing authorization header", "UNAUTHORIZED")
			c.Abort()
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			utils.Fail(c, http.StatusUnauthorized, "invalid authorization header format", "UNAUTHORIZED")
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(jwtSecret, parts[1])
		if err != nil {
			utils.Fail(c, http.StatusUnauthorized, "invalid or expired token", "UNAUTHORIZED")
			c.Abort()
			return
		}

		c.Set(ctxUserID, claims.UserID)
		c.Set(ctxEmail, claims.Email)
		c.Set(ctxRole, claims.Role)
		c.Next()
	}
}

// GetUserID reads the authenticated user's id from context. Only valid
// on routes behind the Auth middleware.
func GetUserID(c *gin.Context) (uuid.UUID, bool) {
	v, ok := c.Get(ctxUserID)
	if !ok {
		return uuid.Nil, false
	}
	id, ok := v.(uuid.UUID)
	return id, ok
}

// GetUserRole reads the authenticated user's role from context.
func GetUserRole(c *gin.Context) string {
	v, ok := c.Get(ctxRole)
	if !ok {
		return ""
	}
	role, _ := v.(string)
	return role
}
