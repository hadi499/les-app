package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleInter, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak: role tidak sesuai"})
			c.Abort()
			return
		}
		
		roleStr, ok := roleInter.(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak: format role tidak valid"})
			c.Abort()
			return
		}
		
		roleStr = strings.ToLower(strings.TrimSpace(roleStr))
		reqRole := strings.ToLower(strings.TrimSpace(requiredRole))
		
		if roleStr != reqRole && roleStr != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak: role tidak sesuai"})
			c.Abort()
			return
		}
		c.Next()
	}
}
