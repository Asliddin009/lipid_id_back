package routes

import (
	"notes/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(noteHandler *handler.Handler) *gin.Engine {
	router := gin.Default()

	noteAPI := router.Group("/notes")
	noteAPI.Use(noteHandler.GetJWTMiddleware())
	{
		// Создание записи
		noteAPI.POST("/note", noteHandler.CreateNote)
		// Получение записи по ID
		noteAPI.GET("/note/:id", noteHandler.GetNoteByID)
		// Редактирование записи
		noteAPI.PUT("/note/:id", noteHandler.UpdateNote)
		// Удаление записи
		noteAPI.DELETE("/note/:id", noteHandler.DeleteNote)
		// Получение списка всех записей пользователя
		noteAPI.GET("/notes", noteHandler.GetAllNotes)
	}

	return router
}
