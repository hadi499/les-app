package main

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ScoreQuiz struct {
	ID       uint
	Username string
	QuizID   uint
	Score    int
}

func main() {
	dbPassword := "admin123" // or what is the password? The default might be postgres or root
	dsn := fmt.Sprintf("host=localhost user=postgres password=%s dbname=lesblg_db port=5432 sslmode=disable TimeZone=Asia/Jakarta", dbPassword)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to db:", err)
	}

	var scores []ScoreQuiz
	db.Find(&scores)
	fmt.Printf("Found %d scores\n", len(scores))
	for _, s := range scores {
		fmt.Printf("Score: %+v\n", s)
	}
}
