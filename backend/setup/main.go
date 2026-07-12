package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

func main() {
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Gagal membaca credentials.json: %v", err)
	}

	config, err := google.ConfigFromJSON(b, drive.DriveScope)
	if err != nil {
		log.Fatalf("Gagal mem-parsing credentials.json: %v", err)
	}

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("\n=== LANGKAH 1: Otorisasi Google Drive ===\n")
	fmt.Printf("Buka link berikut di browser Anda:\n\n%v\n\n", authURL)
	
	fmt.Println("=== LANGKAH 2: Masukkan Kode ===")
	fmt.Println("Setelah login dan menyetujui, Anda mungkin akan melihat pesan error 'localhost refused to connect' di browser.")
	fmt.Println("Itu normal! Lihat URL di address bar browser Anda, copy KODE yang ada setelah '?code=' dan paste di bawah ini.")
	fmt.Print("\nMasukkan KODE di sini: ")

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Gagal membaca kode: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Gagal mendapatkan token: %v", err)
	}

	f, err := os.OpenFile("token.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Gagal membuat token.json: %v", err)
	}
	defer f.Close()

	json.NewEncoder(f).Encode(tok)
	fmt.Println("\n=============================================")
	fmt.Println("BERHASIL! File token.json telah dibuat.")
	fmt.Println("Sekarang Anda bisa menjalankan server dengan 'go run main.go'")
	fmt.Println("=============================================")
}
