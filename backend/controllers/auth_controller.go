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
	"sync"
)

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

var (
	loginAttempts = make(map[string]*AttemptInfo)
	attemptsMu    sync.Mutex
)

type AttemptInfo struct {
	Count        int
	FirstAttempt time.Time
	BlockedUntil time.Time
}

func checkRateLimit(ip string) (bool, time.Duration) {
	attemptsMu.Lock()
	defer attemptsMu.Unlock()

	now := time.Now()
	info, exists := loginAttempts[ip]

	if !exists {
		return true, 0
	}

	if now.Before(info.BlockedUntil) {
		return false, info.BlockedUntil.Sub(now)
	}

	if now.Sub(info.FirstAttempt) > time.Minute {
		info.Count = 0
		info.FirstAttempt = now
		info.BlockedUntil = time.Time{}
	}

	return true, 0
}

func recordFailedAttempt(ip string) {
	attemptsMu.Lock()
	defer attemptsMu.Unlock()

	now := time.Now()
	info, exists := loginAttempts[ip]

	if !exists || now.Sub(info.FirstAttempt) > time.Minute {
		loginAttempts[ip] = &AttemptInfo{
			Count:        1,
			FirstAttempt: now,
		}
		return
	}

	info.Count++
	if info.Count >= 5 {
		info.BlockedUntil = now.Add(10 * time.Minute)
	}
}

func resetAttempts(ip string) {
	attemptsMu.Lock()
	defer attemptsMu.Unlock()
	delete(loginAttempts, ip)
}

func init() {
	// Menjalankan proses pembersihan (garbage collection) di background setiap 10 menit
	go func() {
		for {
			time.Sleep(10 * time.Minute)
			cleanupRateLimiter()
		}
	}()
}

// cleanupRateLimiter menghapus data IP yang sudah tidak relevan untuk mencegah memory leak
func cleanupRateLimiter() {
	attemptsMu.Lock()
	defer attemptsMu.Unlock()

	now := time.Now()
	for ip, info := range loginAttempts {
		// Hapus dari memori jika:
		// 1. Masa blokir sudah habis (atau tidak pernah diblokir) DAN
		// 2. Sudah lewat 1 menit sejak percobaan pertama (window expired)
		if (info.BlockedUntil.IsZero() || now.After(info.BlockedUntil)) && now.Sub(info.FirstAttempt) > time.Minute {
			delete(loginAttempts, ip)
		}
	}
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

	if len(input.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password minimal 6 karakter"})
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

	clientIP := c.ClientIP()
	allowed, waitTime := checkRateLimit(clientIP)
	if !allowed {
		minutes := int(waitTime.Minutes())
		if minutes == 0 {
			minutes = 1
		}
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Terlalu banyak percobaan gagal. Silakan coba lagi setelah " + waitTime.Round(time.Second).String()})
		return
	}

	var user models.User
	if err := database.DB.Where("username = ? ", input.Username).First(&user).Error; err != nil {
		recordFailedAttempt(clientIP)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau password salah"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		recordFailedAttempt(clientIP)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username/email atau password salah"})
		return
	}

	resetAttempts(clientIP)

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

// ChangePassword — ubah password user yang sedang login
func ChangePassword(c *gin.Context) {
	var input struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(input.NewPassword) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password baru minimal 6 karakter"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.OldPassword)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password lama salah"})
		return
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal meng-hash password baru"})
		return
	}

	// Update password
	user.Password = string(hashedPassword)
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan password baru"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password berhasil diubah"})
}
