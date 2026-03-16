package routes

import (
	"auth/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(h *handler.Handler) *gin.Engine {
	router := gin.Default()

	auth := router.Group("/auth")
	{
		// Публичные endpoints (без авторизации)
		auth.POST("/register", h.RegisterUser)
		auth.POST("/login", h.LoginUser)
		auth.POST("/refresh", h.RefreshToken)

		// Защищенные endpoints (требуют авторизации)
		protected := auth.Group("/")
		protected.Use(h.RequireAuth()) // Применяем middleware аутентификации
		{
			protected.GET("/user", h.GetUserInfo)
			protected.PUT("/user", h.UpdateUser)
			protected.DELETE("/user", h.DeleteUser)
		}
	}
	return router

}
