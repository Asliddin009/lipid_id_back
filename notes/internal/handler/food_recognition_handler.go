package handler

import (
    "net/http"
    "notes/internal/errors"
    "notes/internal/models"

    "github.com/gin-gonic/gin"
)

// AnalyzeFood — POST /food-recognition/analyze
// Заглушка: возвращает моковые данные. Интеграция с ML-сервисом будет добавлена позже.
func (h *Handler) AnalyzeFood(c *gin.Context) {
    _, err := h.extractAuthorID(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": errors.MsgMissingUserID})
        return
    }

    // Получаем файл из multipart/form-data
    file, err := c.FormFile("image")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error":   "Изображение обязательно",
            "details": err.Error(),
        })
        return
    }

    // TODO: Отправить файл на ML-сервис для распознавания
    _ = file

    // Пока возвращаем заглушку
    result := models.FoodRecognitionResult{
        DishName:      "Неизвестное блюдо",
        Calories:      0,
        Proteins:      0,
        Fats:          0,
        SaturatedFats: 0,
        Carbs:         0,
        LipidRating:   "unknown",
        RatingImpact:  0,
        Confidence:    0,
    }

    c.JSON(http.StatusOK, result)
}