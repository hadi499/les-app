//go:build ignore

package main

import (
	"backend/config"
	"backend/database"
	"backend/models"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	config.InitConfig()
	database.Connect()
	database.DB.AutoMigrate(&models.User{}, &models.Subject{}, &models.Exam{})

	// Ensure 3 students exist
	var students []models.User
	for i := 1; i <= 3; i++ {
		username := fmt.Sprintf("student_tester_%d", i)
		var user models.User
		if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
			user = models.User{
				Username: username,
				Password: "password123",
				Role:     "student",
			}
			database.DB.Create(&user)
		}
		students = append(students, user)
	}

	// Ensure 2 subjects exist
	subjectNames := []string{"Matematika", "Bahasa Inggris"}
	var subjects []models.Subject
	for _, name := range subjectNames {
		var sub models.Subject
		if err := database.DB.Where("name = ?", name).First(&sub).Error; err != nil {
			sub = models.Subject{Name: name}
			database.DB.Create(&sub)
		}
		subjects = append(subjects, sub)
	}

	// Seed 20 entries per student
	rand.Seed(time.Now().UnixNano())
	now := time.Now()

	for _, student := range students {
		// Clear existing fake exams for this student to avoid explosion on multiple runs
		database.DB.Where("user_id = ?", student.ID).Delete(&models.Exam{})

		for i := 0; i < 20; i++ {
			// random score between 50 and 100
			score := rand.Intn(51) + 50
			// random date within the last 60 days
			daysAgo := rand.Intn(60)
			examDate := now.AddDate(0, 0, -daysAgo)

			// pick random subject
			subject := subjects[rand.Intn(len(subjects))]

			exam := models.Exam{
				UserID:    student.ID,
				SubjectID: subject.ID,
				ExamDate:  examDate,
				ExamName:  fmt.Sprintf("Ujian Harian %d", i+1),
				Score:     score,
			}
			database.DB.Create(&exam)
		}
	}
	fmt.Println("Seeding complete! Added fake students, subjects, and exams.")
}
