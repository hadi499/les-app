package controllers

import (
	"net/http"
	"time"

	"backend/config"
	"backend/database"
	"backend/models/kuisapp"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type KuisAppClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Status   string `json:"status"`
	jwt.RegisteredClaims
}

func RegisterKuisApp(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Status   string `json:"status"` // "pelajar" atau "umum"
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingUser kuisapp.User
	if err := database.DB.Where("username = ?", input.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username sudah digunakan"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengenkripsi password"})
		return
	}

	status := input.Status
	if status == "" {
		status = "umum" // default
	}

	user := kuisapp.User{
		Username: input.Username,
		Password: string(hashedPassword),
		Role:     "user",
		Status:   status,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Registrasi berhasil", "user": user})
}

func LoginKuisApp(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user kuisapp.User
	if err := database.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau password salah"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau password salah"})
		return
	}

	if user.IsSuspended {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akun Anda telah disuspend. Hubungi administrator."})
		return
	}

	expirationTime := time.Now().Add(72 * time.Hour)
	claims := &KuisAppClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		Status:   user.Status,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.JWTKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal generate token"})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		"kuisapp_token",
		tokenString,
		int(72*time.Hour.Seconds()),
		"/",
		"",
		false, 
		true,  
	)

	c.JSON(http.StatusOK, gin.H{"message": "Login berhasil", "user": user})
}

func LogoutKuisApp(c *gin.Context) {
	c.SetCookie("kuisapp_token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logout berhasil"})
}

func GetKuisAppMe(c *gin.Context) {
	tokenString, err := c.Cookie("kuisapp_token")
	if err != nil || tokenString == "" {
		c.JSON(http.StatusOK, gin.H{"authenticated": false})
		return
	}

	claims := &KuisAppClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return config.JWTKey, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusOK, gin.H{"authenticated": false})
		return
	}

	var user kuisapp.User
	if err := database.DB.Where("id = ?", claims.UserID).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"authenticated": false})
		return
	}

	if user.IsSuspended {
		c.JSON(http.StatusOK, gin.H{"authenticated": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"authenticated": true,
		"user": user,
	})
}
