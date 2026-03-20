package routes

import (
	"data/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(noteHandler *handler.Handler) *gin.Engine {
	router := gin.Default()

	// JWT middleware для всех защищённых роутов
	auth := router.Group("")
	auth.Use(noteHandler.GetJWTMiddleware())

	// === Notes ===
	noteAPI := auth.Group("")
	{
		noteAPI.POST("/note", noteHandler.CreateNote)
		noteAPI.GET("/note/:id", noteHandler.GetNoteByID)
		noteAPI.PUT("/note/:id", noteHandler.UpdateNote)
		noteAPI.DELETE("/note/:id", noteHandler.DeleteNote)
		noteAPI.GET("/notes", noteHandler.GetAllNotes)
	}

	// === Food Entries ===
	foodAPI := auth.Group("/food-entries")
	{
		foodAPI.POST("", noteHandler.CreateFoodEntry)
		foodAPI.GET("", noteHandler.GetFoodEntries)
		foodAPI.DELETE("/:id", noteHandler.DeleteFoodEntry)
		foodAPI.GET("/summary", noteHandler.GetDailySummary)
		foodAPI.GET("/summary/weekly", noteHandler.GetWeeklySummary)
	}

	// === Food Recognition ===
	auth.POST("/food-recognition/analyze", noteHandler.AnalyzeFood)

	// === Dashboard ===
	auth.GET("/home/dashboard", noteHandler.GetDashboard)

	// === Charts ===
	chartsAPI := auth.Group("/charts")
	{
		chartsAPI.GET("/lipid-trend", noteHandler.GetLipidTrend)
		chartsAPI.GET("/nutrition-trend", noteHandler.GetNutritionTrend)
		chartsAPI.GET("/weight-trend", noteHandler.GetWeightTrend)
	}

	// === Analyses ===
	analysesAPI := auth.Group("/analyses")
	{
		analysesAPI.POST("", noteHandler.CreateAnalysis)
		analysesAPI.GET("", noteHandler.GetAnalyses)
		analysesAPI.GET("/:id", noteHandler.GetAnalysisByID)
		analysesAPI.DELETE("/:id", noteHandler.DeleteAnalysis)
	}

	// === Weight ===
	weightAPI := auth.Group("/weight")
	{
		weightAPI.POST("", noteHandler.CreateWeightEntry)
		weightAPI.GET("", noteHandler.GetWeightEntries)
		weightAPI.DELETE("/:id", noteHandler.DeleteWeightEntry)
	}

	// === Devices & Notifications ===
	auth.POST("/devices", noteHandler.RegisterDevice)
	notifAPI := auth.Group("/notifications")
	{
		notifAPI.GET("", noteHandler.GetNotifications)
		notifAPI.PATCH("/:id/read", noteHandler.MarkNotificationRead)
	}

	return router
}
