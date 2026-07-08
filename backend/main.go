package main

import (
	"log"

	"backend/config"
	"backend/database"
	"backend/models"
	"backend/routes"

	"os"
	"strings"
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
	database.DB.AutoMigrate(&models.User{}, &models.BlacklistedToken{}, &models.LessonProgress{}, &models.GameHighScore{}, &models.GameHistory{}, &models.LessonHistory{}, &models.CardFolder{}, &models.Card{}, &models.Exam{}, &models.Subject{}, &models.Quiz{}, &models.Question{}, &models.ScoreQuiz{}, &models.Folder{}, &models.Note{}, &models.Absence{}, &models.TodoList{}, &models.TodoItem{}, &models.SystemSetting{}, &models.WritingProgress{})

	// Seeding pengaturan awal
	var settingPaud models.SystemSetting
	if err := database.DB.Where("key = ?", "is_class_open_paud").First(&settingPaud).Error; err != nil {
		database.DB.Create(&models.SystemSetting{Key: "is_class_open_paud", Value: "true"})
	}

	var settingSd models.SystemSetting
	if err := database.DB.Where("key = ?", "is_class_open_sd").First(&settingSd).Error; err != nil {
		database.DB.Create(&models.SystemSetting{Key: "is_class_open_sd", Value: "true"})
	}

	// Background job untuk membersihkan token blacklist yang kedaluwarsa setiap 1 jam
	go func() {
		for {
			database.DB.Where("expires_at < ?", time.Now()).Delete(&models.BlacklistedToken{})
			time.Sleep(1 * time.Hour)
		}
	}()

	// Background job untuk menghapus gambar bukti nilai ujian yang usianya lebih dari 1 bulan
	go func() {
		for {
			var exams []models.Exam
			// Batas waktu penghapusan: 30 hari (1 bulan) sejak update terakhir
			batasWaktu := time.Now().AddDate(0, -1, 0)
			
			// Cari data ujian yang memiliki gambar dan diupdate sebelum batas waktu
			database.DB.Where("image != ? AND image IS NOT NULL AND updated_at < ?", "", batasWaktu).Find(&exams)
			
			for _, exam := range exams {
				if exam.Image != "" {
					// Hapus slash di awal agar path menjadi relatif (misal: /uploads/file.jpg -> uploads/file.jpg)
					filePath := strings.TrimPrefix(exam.Image, "/")
					_ = os.Remove(filePath)
					
					// Hapus path gambar dari database tanpa menghapus data nilainya
					database.DB.Model(&exam).Update("image", "")
				}
			}
			// Cek setiap 24 jam sekali
			time.Sleep(24 * time.Hour)
		}
	}()

	// Background job untuk menghapus gambar bukti perkembangan menulis yang usianya lebih dari 1 bulan
	go func() {
		for {
			var progresses []models.WritingProgress
			batasWaktu := time.Now().AddDate(0, -1, 0)
			
			database.DB.Where("image != ? AND image IS NOT NULL AND updated_at < ?", "", batasWaktu).Find(&progresses)
			
			for _, wp := range progresses {
				if wp.Image != "" {
					filePath := strings.TrimPrefix(wp.Image, "/")
					_ = os.Remove(filePath)
					database.DB.Model(&wp).Update("image", "")
				}
			}
			time.Sleep(24 * time.Hour)
		}
	}()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://192.168.18.2:5173", "http://127.0.0.1:5173", "http://172.27.210.181:5173", "http://43.129.51.155", "https://lesbalonggarut.my.id", "https://www.lesbalonggarut.my.id"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Setup semua route
	routes.SetupRoutes(r)

	// Menyediakan akses publik ke folder uploads
	r.Static("/uploads", "./uploads")

	log.Println("Server berjalan di port 8080...")
	r.Run(":8080")
}
