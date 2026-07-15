package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/api/auth/register", middleware.AuthMiddleware(), middleware.RoleMiddleware("teacher"), controllers.Register)
	r.POST("/api/auth/login", controllers.Login)
	r.POST("/api/auth/logout", middleware.AuthMiddleware(), controllers.Logout)
	r.PUT("/api/auth/change-password", middleware.AuthMiddleware(), controllers.ChangePassword)

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

	// Exams API routes
	exams := r.Group("/api/exams")
	exams.Use(middleware.AuthMiddleware())
	{
		exams.GET("", controllers.GetExams)
		exams.GET("/:id", controllers.GetExamByID)

		teacherExams := exams.Group("")
		teacherExams.Use(middleware.RoleMiddleware("teacher"))
		{
			teacherExams.POST("", controllers.CreateExam)
			teacherExams.POST("/bulk-delete", controllers.BulkDeleteExams)
			teacherExams.PUT("/:id", controllers.UpdateExam)
			teacherExams.DELETE("/:id", controllers.DeleteExam)
		}
	}

	// Subjects API routes
	subjects := r.Group("/api/subjects")
	subjects.Use(middleware.AuthMiddleware())
	{
		subjects.GET("", controllers.GetSubjects)

		teacherSubjects := subjects.Group("")
		teacherSubjects.Use(middleware.RoleMiddleware("teacher"))
		{
			teacherSubjects.POST("", controllers.CreateSubject)
			teacherSubjects.PUT("/:id", controllers.UpdateSubject)
			teacherSubjects.DELETE("/:id", controllers.DeleteSubject)
		}
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

	// Writing Progress API routes
	writingProgress := r.Group("/api/writing-progress")
	writingProgress.Use(middleware.AuthMiddleware())
	{
		writingProgress.GET("", controllers.GetWritingProgresses)
		writingProgress.GET("/:id", controllers.GetWritingProgressByID)

		teacherWritingProgress := writingProgress.Group("")
		teacherWritingProgress.Use(middleware.RoleMiddleware("teacher"))
		{
			teacherWritingProgress.POST("", controllers.CreateWritingProgress)
			teacherWritingProgress.POST("/backup", controllers.BackupToDrive)
			teacherWritingProgress.POST("/bulk-delete", controllers.BulkDeleteWritingProgress)
			teacherWritingProgress.PUT("/:id", controllers.UpdateWritingProgress)
			teacherWritingProgress.DELETE("/:id", controllers.DeleteWritingProgress)
		}
	}

	// Card Folders API routes (Teacher only)
	cardFolders := r.Group("/api/card-folders")
	cardFolders.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("teacher"))
	{
		cardFolders.GET("", controllers.GetCardFolders)
		cardFolders.POST("", controllers.CreateCardFolder)
		cardFolders.PUT("/:id", controllers.UpdateCardFolder)
		cardFolders.DELETE("/:id", controllers.DeleteCardFolder)
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

	// Absences API routes
	absences := r.Group("/api/absences")
	absences.Use(middleware.AuthMiddleware())
	{
		absences.POST("", middleware.RoleMiddleware("teacher"), controllers.CreateAbsence)
		absences.GET("/recap", controllers.GetAbsenceRecap)
		absences.GET("/user/:id", controllers.GetAbsenceHistory)
		absences.POST("/reset", middleware.RoleMiddleware("teacher"), controllers.ResetAbsences)
		absences.PUT("/:id", middleware.RoleMiddleware("teacher"), controllers.UpdateAbsence)
		absences.DELETE("/:id", middleware.RoleMiddleware("teacher"), controllers.DeleteAbsence)
	}

	// Todolist API routes
	todolists := r.Group("/api/todolists")
	todolists.Use(middleware.AuthMiddleware())
	{
		todolists.GET("", controllers.GetTodoLists)
		todolists.GET("/:id", controllers.GetTodoList)

		teacherTodos := todolists.Group("")
		teacherTodos.Use(middleware.RoleMiddleware("teacher"))
		{
			teacherTodos.POST("", controllers.CreateTodoList)
			teacherTodos.PUT("/:id", controllers.UpdateTodoList)
			teacherTodos.DELETE("/:id", controllers.DeleteTodoList)
			teacherTodos.POST("/:id/items", controllers.CreateTodoItem)
			teacherTodos.PUT("/:id/items/:item_id", controllers.ToggleTodoItem)
			teacherTodos.DELETE("/:id/items/:item_id", controllers.DeleteTodoItem)
		}
	}

	// Scores API routes (User submission)
	scores := r.Group("/api/scores")
	scores.Use(middleware.AuthMiddleware())
	{
		scores.POST("/quizzes", controllers.SubmitQuizScore)
		scores.GET("/quizzes", controllers.GetMyQuizScores)
	}

	// System API routes (Teacher only)
	system := r.Group("/api/system")
	system.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("teacher"))
	{
		system.GET("/info", controllers.GetSystemInfo)
	}

	// Logs API routes (Teacher only)
	logsAPI := r.Group("/api/logs")
	logsAPI.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("teacher"))
	{
		logsAPI.GET("", controllers.GetLogs)
	}

	// Settings API routes
	r.GET("/api/settings", controllers.GetSettings) // Publik
	
	settingsAdmin := r.Group("/api/settings")
	settingsAdmin.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("teacher"))
	{
		settingsAdmin.PUT("/:key", controllers.UpdateSetting)
	}

	// Chat API routes
	chat := r.Group("/api/chat")
	chat.Use(middleware.AuthMiddleware())
	{
		chat.GET("/contacts", controllers.GetContacts)
		chat.GET("/history/:userId", controllers.GetChatHistory)
		chat.GET("/unread-count", controllers.GetUnreadCount)
		chat.DELETE("/messages/:id", controllers.DeleteMessage)
	}
	r.GET("/ws/chat", middleware.AuthMiddleware(), controllers.HandleChatWebSocket)
}

