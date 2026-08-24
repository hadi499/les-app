package middleware

import (
	"net/http"
	"strings"

	"backend/config"
	"backend/database"
	"backend/models/kuisapp"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type KuisAppClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Status   string `json:"status"`
	jwt.RegisteredClaims
}

func KuisAppAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("kuisapp_token")
		
		if err != nil || tokenString == "" {
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token kuisapp tidak ditemukan"})
			c.Abort()
			return
		}

		claims := &KuisAppClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return config.JWTKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token kuisapp tidak valid atau sudah kadaluarsa"})
			c.Abort()
			return
		}

		var user kuisapp.User
		if err := database.DB.First(&user, claims.UserID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak ditemukan"})
			c.Abort()
			return
		}

		if user.IsSuspended {
			c.JSON(http.StatusForbidden, gin.H{"error": "Akun Anda telah disuspend. Hubungi administrator."})
			c.Abort()
			return
		}

		c.Set("kuisapp_user", user)
		c.Next()
	}
}

func KuisAppRoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userObj, exists := c.Get("kuisapp_user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Tidak terautentikasi"})
			c.Abort()
			return
		}

		user := userObj.(kuisapp.User)
		if user.Role != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak: role tidak sesuai"})
			c.Abort()
			return
		}
		c.Next()
	}
}
