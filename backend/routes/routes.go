package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/api/auth/register", middleware.AuthMiddleware(), middleware.RoleMiddleware("teacher"), controllers.Register)
	r.POST("/api/auth/public-register", controllers.PublicRegister)
	r.POST("/api/auth/login", controllers.Login)
	r.POST("/api/auth/logout", middleware.AuthMiddleware(), controllers.Logout)
	r.PUT("/api/auth/change-password", middleware.AuthMiddleware(), controllers.ChangePassword)

	// Kuis App Auth Routes
	r.POST("/api/kuisapp/register", controllers.RegisterKuisApp)
	r.POST("/api/kuisapp/login", controllers.LoginKuisApp)
	r.POST("/api/kuisapp/logout", middleware.KuisAppAuthMiddleware(), controllers.LogoutKuisApp)
	r.GET("/api/kuisapp/me", controllers.GetKuisAppMe)

	// Kuis App API Routes
	kuisappApi := r.Group("/api/kuisapp")
	kuisappApi.Use(middleware.KuisAppAuthMiddleware())
	{
		kuisappApi.POST("/change-password", controllers.ChangePasswordKuisApp)
		// Categories
		kuisappApi.GET("/categories", controllers.GetKuisAppCategories)
		// Quizzes
		kuisappApi.GET("/quizzes", controllers.GetKuisAppQuizzes)
		kuisappApi.GET("/quizzes/:id", controllers.GetKuisAppQuiz)
		// Results (User)
		kuisappApi.POST("/quizzes/:id/submit", controllers.SubmitKuisAppQuiz)
		kuisappApi.GET("/my-results", controllers.GetKuisAppMyResults)

		// Admin only routes
		adminApi := kuisappApi.Group("")
		adminApi.Use(middleware.KuisAppRoleMiddleware("admin"))
		{
			adminApi.POST("/categories", controllers.CreateKuisAppCategory)
			adminApi.PUT("/categories/:id", controllers.UpdateKuisAppCategory)
			adminApi.DELETE("/categories/:id", controllers.DeleteKuisAppCategory)
			
			adminApi.POST("/upload", controllers.UploadKuisAppImage)
			
			adminApi.POST("/quizzes", controllers.CreateKuisAppQuiz)
			adminApi.PUT("/quizzes/:id", controllers.UpdateKuisAppQuiz)
			adminApi.DELETE("/quizzes/:id", controllers.DeleteKuisAppQuiz)
			adminApi.POST("/quizzes/:id/duplicate", controllers.DuplicateKuisAppQuiz)
			
			adminApi.PUT("/quizzes/:id/questions/bulk", controllers.BulkSaveKuisAppQuestions)
			
			adminApi.POST("/quizzes/:id/questions", controllers.CreateKuisAppQuestion)
			adminApi.PUT("/questions/:question_id", controllers.UpdateKuisAppQuestion)
			adminApi.DELETE("/questions/:question_id", controllers.DeleteKuisAppQuestion)
			
			adminApi.GET("/all-results", controllers.GetKuisAppAllResults)
			
			// User Management
			adminApi.GET("/users", controllers.GetKuisAppUsers)
			adminApi.POST("/users", controllers.CreateKuisAppUser)
			adminApi.PUT("/users/:id", controllers.UpdateKuisAppUser)
			adminApi.DELETE("/users/:id", controllers.DeleteKuisAppUser)
			adminApi.POST("/users/:id/reset-password", controllers.ResetKuisAppUserPassword)
			adminApi.POST("/users/:id/reset-points", controllers.ResetKuisAppUserPoints)
			adminApi.POST("/reset-all-points-and-history", controllers.ResetAllKuisAppPointsAndHistory)
			adminApi.PUT("/users/:id/suspend", controllers.ToggleKuisAppUserSuspend)
		}
	}

	// Route /me tidak menggunakan AuthMiddleware agar bisa mereturn 200 dengan status authenticated = false
	// alih-alih mereturn 401 Unauthorized yang akan memicu log merah di browser.
	r.GET("/me", controllers.Me)

	// Users API routes (Teacher only)
	users := r.Group("/api/users")
	users.Use(middleware.AuthMiddleware())
	{
		users.GET("", controllers.GetUsers)
		users.GET("/:id", controllers.GetUserByID)
		users.PUT("/:id", controllers.UpdateUser)
		users.DELETE("/:id", controllers.DeleteUser)
		users.POST("/:id/reset-password", controllers.ResetUserPassword)
		users.POST("/reset-points", middleware.RoleMiddleware("teacher"), controllers.ResetAllPoints)
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
		teacherOnly.Use(middleware.RoleMiddleware("admin"))
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
		uploads.DELETE("", controllers.DeleteUploadedImage)
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

		manageSubjects := subjects.Group("")
		manageSubjects.Use(middleware.RoleMiddleware("teacher"))
		{
			manageSubjects.POST("", controllers.CreateSubject)
			manageSubjects.PUT("/:id", controllers.UpdateSubject)
			manageSubjects.DELETE("/:id", controllers.DeleteSubject)
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
			teacherQuizzes.DELETE("/:id/scores", controllers.ResetQuizScores)
			teacherQuizzes.POST("/:id/duplicate", controllers.DuplicateQuiz)
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
			teacherWritingProgress.POST("/bulk-delete", controllers.BulkDeleteWritingProgress)
			teacherWritingProgress.PUT("/:id", controllers.UpdateWritingProgress)
			teacherWritingProgress.DELETE("/:id", controllers.DeleteWritingProgress)
		}

		adminWritingProgress := writingProgress.Group("")
		adminWritingProgress.Use(middleware.RoleMiddleware("admin"))
		{
			adminWritingProgress.POST("/backup", controllers.BackupToDrive)
		}
	}

	// Card Folders API routes (Teacher only)
	cardFolders := r.Group("/api/card-folders")
	cardFolders.Use(middleware.AuthMiddleware())
	{
		cardFolders.GET("", controllers.GetCardFolders)

		adminFolders := cardFolders.Group("")
		adminFolders.Use(middleware.RoleMiddleware("admin"))
		{
			adminFolders.POST("", controllers.CreateCardFolder)
			adminFolders.PUT("/:id", controllers.UpdateCardFolder)
			adminFolders.DELETE("/:id", controllers.DeleteCardFolder)
		}
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

	// Materi API routes
	materis := r.Group("/api/materis")
	materis.Use(middleware.AuthMiddleware())
	{
		materis.GET("", controllers.GetMateris)
		materis.GET("/:id", controllers.GetMateriByID)

		teacherMateris := materis.Group("")
		teacherMateris.Use(middleware.RoleMiddleware("teacher"))
		{
			teacherMateris.POST("", controllers.CreateMateri)
			teacherMateris.PUT("/:id", controllers.UpdateMateri)
			teacherMateris.DELETE("/:id", controllers.DeleteMateri)
		}
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

		todolists.POST("", controllers.CreateTodoList)
		todolists.PUT("/:id", controllers.UpdateTodoList)
		todolists.DELETE("/:id", controllers.DeleteTodoList)
		todolists.POST("/:id/items", controllers.CreateTodoItem)
		todolists.PUT("/:id/items/:item_id", controllers.ToggleTodoItem)
		todolists.DELETE("/:id/items/:item_id", controllers.DeleteTodoItem)
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
	system.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"))
	{
		system.GET("/info", controllers.GetSystemInfo)
	}

	// Logs API routes (Teacher only)
	logsAPI := r.Group("/api/logs")
	logsAPI.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"))
	{
		logsAPI.GET("", controllers.GetLogs)
	}

	// Quotes API routes
	quotes := r.Group("/api/quotes")
	{
		quotes.GET("", controllers.GetPublicQuotes) // public
		
		adminQuotes := quotes.Group("")
		adminQuotes.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"))
		{
			adminQuotes.GET("/all", controllers.GetAllQuotes)
			adminQuotes.POST("", controllers.CreateQuote)
			adminQuotes.PUT("/:id", controllers.UpdateQuote)
			adminQuotes.DELETE("/:id", controllers.DeleteQuote)
		}
	}

	// Settings API routes
	r.GET("/api/settings", controllers.GetSettings) // Publik
	
	settingsAdmin := r.Group("/api/settings")
	settingsAdmin.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("teacher"))
	{
		settingsAdmin.PUT("/:key", controllers.UpdateSetting)
	}

	// Parents API routes
	parents := r.Group("/api/parents")
	parents.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("parent"))
	{
		parents.GET("/my-children", controllers.GetMyChildren)
	}

	// Chat API routes
	chat := r.Group("/api/chat")
	chat.Use(middleware.AuthMiddleware())
	{
		chat.GET("/contacts", controllers.GetContacts)
		chat.GET("/history/:userId", controllers.GetChatHistory)
		chat.GET("/unread-count", controllers.GetUnreadCount)
		chat.DELETE("/messages/:id", controllers.DeleteMessage)
		chat.POST("/broadcast", middleware.RoleMiddleware("teacher"), controllers.BroadcastMessage)
	}
	r.GET("/ws/chat", middleware.AuthMiddleware(), controllers.HandleChatWebSocket)
}

