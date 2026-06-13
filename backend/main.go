package main

import (
	"log"

	"backend/config"
	"backend/database"
	"backend/models"
	"backend/routes"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Muat file .env sekali di awal aplikasi
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// 2. Inisialisasi konfigurasi global
	config.InitConfig()

	// 3. Koneksi ke Database
	database.Connect()

	// Migrasi otomatis untuk memastikan tabel ada
	database.DB.AutoMigrate(&models.User{}, &models.BlacklistedToken{}, &models.LessonProgress{}, &models.GameHighScore{}, &models.GameHistory{}, &models.LessonHistory{}, &models.Card{}, &models.Exam{}, &models.Subject{})

	// Background job untuk membersihkan token blacklist yang kedaluwarsa setiap 1 jam
	go func() {
		for {
			database.DB.Where("expires_at < ?", time.Now()).Delete(&models.BlacklistedToken{})
			time.Sleep(1 * time.Hour)
		}
	}()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://192.168.18.2:5173"}, // Alamat Svelte Anda
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // Wajib agar browser mengirim cookie
		MaxAge:           12 * time.Hour,
	}))

	// Setup semua route
	routes.SetupRoutes(r)

	// Menyediakan akses publik ke folder uploads
	r.Static("/uploads", "./uploads")

	log.Println("Server berjalan di port 8080...")
	r.Run(":8080")
}
