package middleware

import (
	"backend/config"
	"backend/database"
	"backend/models"
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Baca token dari HttpOnly cookie (bukan Authorization header)
		tokenString, err := c.Cookie("auth_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication cookie tidak ditemukan"})
			c.Abort()
			return
		}

		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return config.JWTKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		hash := sha256.Sum256([]byte(tokenString))
		tokenHash := hex.EncodeToString(hash[:])

		var blacklisted models.BlacklistedToken
		if err := database.DB.Where("token_hash = ?", tokenHash).First(&blacklisted).Error; err == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has been revoked"})
			c.Abort()
			return
		}

		// Fetch user to ensure role is up to date (and handle old tokens)
		role := claims.Role
		var user models.User
		if err := database.DB.Where("id = ?", claims.UserID).First(&user).Error; err == nil {
			role = user.Role
		}

		// simpan ke context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", role)

		c.Next()
	}
}
