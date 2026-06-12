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
}
