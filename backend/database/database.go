package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	// Ambil hanya nilai password dari .env
	dbPassword := os.Getenv("DB_PASSWORD")

	// 3. Susun string DSN, gunakan %s untuk tempat masuknya dbPassword
	dsn := fmt.Sprintf("host=localhost user=hadi password=%s dbname=les_db port=5432 sslmode=disable TimeZone=Asia/Jakarta", dbPassword)

	// 4. Lakukan koneksi ke database
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal koneksi ke database:", err)
	}

	// Setup Connection Pooling untuk Server RAM 2GB
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Gagal mengambil instance sql.DB:", err)
	}

	// SetMaxIdleConns: jumlah maksimal koneksi nganggur (idle) di dalam pool.
	// Nilai 5-10 sudah cukup agar saat traffic sepi, aplikasi tidak menahan banyak RAM PostgreSQL.
	sqlDB.SetMaxIdleConns(5)

	// SetMaxOpenConns: batas maksimal koneksi aktif bersamaan.
	// Server 2GB (dengan OS dan service lain) bisa menahan sekitar 20-40 koneksi PostgreSQL dengan aman.
	// Jika ada 1000 request masuk, Golang akan antrikan otomatis ke dalam 25 koneksi ini.
	sqlDB.SetMaxOpenConns(25)

	// SetConnMaxLifetime: umur maksimal sebuah koneksi.
	// Memaksa koneksi yang terlalu lama hidup untuk direfresh, mencegah masalah stale connection & memory leak ringan.
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Koneksi ke database berhasil dengan Connection Pooling!")
}
