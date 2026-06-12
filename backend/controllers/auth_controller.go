package controllers

import (
	"backend/config"
	"backend/database"
	"backend/models"
	"crypto/sha256"
	"encoding/hex"

	"net/http"

	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// Register — mendaftarkan user baru
func Register(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Role     string `json:"role"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal meng-hash password"})
		return
	}

	// Default role jika tidak diisi
	role := input.Role
	if role == "" {
		role = "student"
	}

	user := models.User{
		Username: input.Username,
		Password: string(hashedPassword),
		Role:     role,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		// Tangani duplicate unique constraint
		if strings.Contains(err.Error(), "duplicate") {
			c.JSON(http.StatusConflict, gin.H{"error": "Username atau email sudah digunakan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendaftarkan user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pendaftaran berhasil",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}

// Login — autentikasi user dan set cookie HttpOnly
func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("username = ? ", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau password salah"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username/email atau password salah"})
		return
	}

	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role, // Mengambil role asli dari database
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)), // 3 hari
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.JWTKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
		return
	}

	// Set SameSite parameter untuk proteksi CSRF
	c.SetSameSite(http.SameSiteLaxMode)

	// Set HttpOnly cookie
	c.SetCookie(
		"auth_token",
		tokenString,
		72*3600, // 3 hari dalam detik
		"/",
		"",
		false, // Secure: true hanya di production (HTTPS)
		true,  // HttpOnly
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}

// Logout — blacklist token dan hapus cookie
func Logout(c *gin.Context) {
	tokenString, err := c.Cookie("auth_token")
	if err == nil && tokenString != "" {
		// Hash token untuk disimpan ke blacklist
		hash := sha256.Sum256([]byte(tokenString))
		tokenHash := hex.EncodeToString(hash[:])

		blacklisted := models.BlacklistedToken{
			TokenHash: tokenHash,
			ExpiresAt: time.Now().Add(72 * time.Hour),
		}
		database.DB.Create(&blacklisted)
	}

	// Hapus cookie dengan SameSiteLaxMode
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("auth_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Logout berhasil"})
}

// Me — ambil data user yang sedang login
func Me(c *gin.Context) {
	tokenString, err := c.Cookie("auth_token")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"authenticated": false})
		return
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return config.JWTKey, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusOK, gin.H{"authenticated": false})
		return
	}

	hash := sha256.Sum256([]byte(tokenString))
	tokenHash := hex.EncodeToString(hash[:])

	var blacklisted models.BlacklistedToken
	if err := database.DB.Where("token_hash = ?", tokenHash).First(&blacklisted).Error; err == nil {
		c.JSON(http.StatusOK, gin.H{"authenticated": false})
		return
	}

	var user models.User
	if err := database.DB.Where("id = ?", claims.UserID).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"authenticated": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"authenticated": true,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}
