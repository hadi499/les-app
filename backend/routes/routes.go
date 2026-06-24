package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/api/auth/register", controllers.Register)
	r.POST("/api/auth/login", controllers.Login)
	r.POST("/api/auth/logout", middleware.AuthMiddleware(), controllers.Logout)
	
	// Route /me tidak menggunakan AuthMiddleware agar bisa mereturn 200 dengan status authenticated = false
	// alih-alih mereturn 401 Unauthorized yang akan memicu log merah di browser.
	r.GET("/me", controllers.Me)

	// Users API routes (Teacher only)
	users := r.Group("/api/users")
	users.Use(middleware.AuthMiddleware())
	{
		users.GET("", controllers.GetUsers)
		users.DELETE("/:id", controllers.DeleteUser)
	}

	// Typing API routes
	typing := r.Group("/api/typing")
	typing.Use(middleware.AuthMiddleware())
	{
		typing.GET("/progress", controllers.GetProgress)
		typing.POST("/progress", controllers.SaveProgress)
		typing.GET("/game-scores", controllers.GetGameScores)
		typing.POST("/game-scores", controllers.SaveGameScore)
		typing.GET("/history/game", controllers.GetGameHistory)
		typing.GET("/history/lesson", controllers.GetLessonHistory)

		// Admin/Teacher routes
		admin := typing.Group("/admin")
		admin.Use(middleware.RoleMiddleware("teacher"))
		{
			admin.GET("/progress", controllers.GetAllLessonProgress)
			admin.GET("/game-scores", controllers.GetAllGameScores)
			admin.GET("/history/game", controllers.GetAllGameHistory)
			admin.GET("/history/lesson", controllers.GetAllLessonHistory)
		}
	}

	// Cards API routes
	cards := r.Group("/api/cards")
	cards.Use(middleware.AuthMiddleware())
	{
		// Read can be done by any authenticated user (e.g. students or teachers viewing cards)
		cards.GET("", controllers.GetCards)
		
		// Teacher only for modifications
		teacherOnly := cards.Group("")
		teacherOnly.Use(middleware.RoleMiddleware("teacher"))
		{
			teacherOnly.POST("", controllers.CreateCard)
			teacherOnly.PUT("/:id", controllers.UpdateCard)
			teacherOnly.DELETE("/:id", controllers.DeleteCard)
			
			teacherOnly.GET("/trash", controllers.GetTrashCards)
			teacherOnly.POST("/trash/:id/restore", controllers.RestoreCard)
			teacherOnly.DELETE("/trash/:id/force", controllers.ForceDeleteCard)
			teacherOnly.DELETE("/trash/empty", controllers.EmptyTrash)
		}
	}

	// Upload API routes (Teacher only)
	uploads := r.Group("/api/upload")
	uploads.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("teacher"))
	{
		uploads.POST("", controllers.UploadImage)
	}

	imagesAPI := r.Group("/api/images")
	imagesAPI.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("teacher"))
	{
		imagesAPI.GET("", controllers.ListImages)
	}

	// Exams API routes (Teacher only)
	exams := r.Group("/api/exams")
	exams.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("teacher"))
	{
		exams.GET("", controllers.GetExams)
		exams.POST("", controllers.CreateExam)
		exams.PUT("/:id", controllers.UpdateExam)
	exams.DELETE("/:id", controllers.DeleteExam)
	}

	// Subjects API routes (Teacher only)
	subjects := r.Group("/api/subjects")
	subjects.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("teacher"))
	{
		subjects.GET("", controllers.GetSubjects)
		subjects.POST("", controllers.CreateSubject)
		subjects.PUT("/:id", controllers.UpdateSubject)
		subjects.DELETE("/:id", controllers.DeleteSubject)
	}

	// Quizzes API routes
	quizzes := r.Group("/api/quizzes")
	quizzes.Use(middleware.AuthMiddleware())
	{
		quizzes.GET("", controllers.GetQuizzes)
		quizzes.GET("/:id", controllers.GetQuizByID)

		teacherQuizzes := quizzes.Group("")
		teacherQuizzes.Use(middleware.RoleMiddleware("teacher"))
		{
			teacherQuizzes.POST("", controllers.CreateQuiz)
			teacherQuizzes.PUT("/:id", controllers.UpdateQuiz)
			teacherQuizzes.DELETE("/:id", controllers.DeleteQuiz)
			teacherQuizzes.GET("/scores", controllers.GetQuizScores)
		}
	}

	// Folders API routes (Teacher only)
	folders := r.Group("/api/folders")
	folders.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("teacher"))
	{
		folders.GET("", controllers.GetFolders)
		folders.POST("", controllers.CreateFolder)
		folders.PUT("/:id", controllers.UpdateFolder)
		folders.DELETE("/:id", controllers.DeleteFolder)
	}

	// Notes API routes (Teacher only)
	notes := r.Group("/api/notes")
	notes.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("teacher"))
	{
		notes.GET("", controllers.GetNotes)
		notes.POST("", controllers.CreateNote)
		notes.PUT("/:id", controllers.UpdateNote)
		notes.DELETE("/:id", controllers.DeleteNote)
	}

	// Absences API routes (Teacher only)
	absences := r.Group("/api/absences")
	absences.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("teacher"))
	{
		absences.POST("", controllers.CreateAbsence)
		absences.GET("/recap", controllers.GetAbsenceRecap)
		absences.GET("/user/:id", controllers.GetAbsenceHistory)
		absences.POST("/reset", controllers.ResetAbsences)
		absences.PUT("/:id", controllers.UpdateAbsence)
		absences.DELETE("/:id", controllers.DeleteAbsence)
	}

	// Scores API routes (User submission)
	scores := r.Group("/api/scores")
	scores.Use(middleware.AuthMiddleware())
	{
		scores.POST("/quizzes", controllers.SubmitQuizScore)
		scores.GET("/quizzes", controllers.GetMyQuizScores)
	}
}
