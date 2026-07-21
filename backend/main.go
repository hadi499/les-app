package main

import (
	"log"

	"backend/config"
	"backend/database"
	"backend/models"
	"backend/routes"
	"backend/services"


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

	// 4. Inisialisasi Google Drive Service (Opsional jika credentials ada)
	if os.Getenv("DRIVE_CREDENTIALS_FILE") != "" {
		if err := services.InitDriveService(os.Getenv("DRIVE_CREDENTIALS_FILE")); err != nil {
			log.Printf("Gagal inisialisasi Google Drive API: %v\n", err)
		} else {
			log.Println("Google Drive API berhasil diinisialisasi")
		}
	}

	// Migrasi otomatis untuk memastikan tabel ada
	database.DB.AutoMigrate(&models.User{}, &models.BlacklistedToken{}, &models.LessonProgress{}, &models.GameHighScore{}, &models.GameHistory{}, &models.LessonHistory{}, &models.CardFolder{}, &models.Card{}, &models.Exam{}, &models.Subject{}, &models.Quiz{}, &models.Question{}, &models.ScoreQuiz{}, &models.Folder{}, &models.Note{}, &models.Absence{}, &models.TodoList{}, &models.TodoItem{}, &models.SystemSetting{}, &models.WritingProgress{}, &models.ChatMessage{}, &models.UserLog{}, &models.Materi{}, &models.Quote{})

	// Seeding pengaturan awal
	var settingPaud models.SystemSetting
	if err := database.DB.Where("key = ?", "is_class_open_paud").First(&settingPaud).Error; err != nil {
		database.DB.Create(&models.SystemSetting{Key: "is_class_open_paud", Value: "true"})
	}

	var settingSd models.SystemSetting
	if err := database.DB.Where("key = ?", "is_class_open_sd").First(&settingSd).Error; err != nil {
		database.DB.Create(&models.SystemSetting{Key: "is_class_open_sd", Value: "true"})
	}

	// Seeding quotes awal
	var countQuotes int64
	database.DB.Model(&models.Quote{}).Count(&countQuotes)
	if countQuotes == 0 {
		initialQuotes := []models.Quote{
			{
				Quote:  "The beautiful thing about learning is that no one can take it away from you.",
				Arti:   "Hal yang indah tentang belajar adalah tidak ada yang bisa mengambilnya darimu.",
				Author: "B.B. King",
			},
			{
				Quote:  "Education is the most powerful weapon which you can use to change the world.",
				Arti:   "Pendidikan adalah senjata paling ampuh yang bisa Anda gunakan untuk mengubah dunia.",
				Author: "Nelson Mandela",
			},
			{
				Quote:  "Success is no accident. It is hard work, perseverance, learning, studying, sacrifice and most of all, love of what you are doing or learning to do.",
				Arti:   "Kesuksesan bukanlah sebuah kebetulan. Ia adalah kerja keras, ketekunan, pembelajaran, pengorbanan, dan yang terpenting, cinta akan apa yang Anda lakukan atau pelajari.",
				Author: "Pelé",
			},
		}
		database.DB.Create(&initialQuotes)
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

	// Background job untuk menghapus data dan gambar bukti perkembangan menulis yang usianya lebih dari 1 bulan
	go func() {
		for {
			var progresses []models.WritingProgress
			// Batas waktu penghapusan: 30 hari yang lalu
			batasWaktu := time.Now().AddDate(0, -1, 0)
			
			// Cari semua data perkembangan menulis yang diupdate sebelum batas waktu
			database.DB.Where("updated_at < ?", batasWaktu).Find(&progresses)
			
			for _, wp := range progresses {
				if wp.Image != "" {
					filePath := strings.TrimPrefix(wp.Image, "/")
					_ = os.Remove(filePath)
				}
				// Hapus record dari database
				database.DB.Delete(&wp)
			}
			// Cek setiap 24 jam sekali
			time.Sleep(24 * time.Hour)
		}
	}()

	// Background job untuk mereset poin semua user setiap awal bulan
	go func() {
		for {
			now := time.Now()
			currentMonth := now.Format("2006-01") // Format YYYY-MM
			
			var setting models.SystemSetting
			err := database.DB.Where("key = ?", "last_point_reset").First(&setting).Error
			
			if err != nil {
				// Belum pernah diset, set ke bulan ini agar tidak mereset bulan berjalan saat fitur ini pertama kali dirilis
				database.DB.Create(&models.SystemSetting{Key: "last_point_reset", Value: currentMonth})
			} else if setting.Value != currentMonth {
				// Bulan sudah berganti, reset semua poin ke 0
				database.DB.Model(&models.User{}).Where("1 = 1").Update("points", 0)
				// Reset semua riwayat nilai kuis (menghapus semua skor agar bisa dikerjakan ulang dari nol)
				database.DB.Where("1 = 1").Delete(&models.ScoreQuiz{})
				// Update bulan terakhir reset
				database.DB.Model(&setting).Update("value", currentMonth)
			}

			// Cek setiap 24 jam sekali
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
