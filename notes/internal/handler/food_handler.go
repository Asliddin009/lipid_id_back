package handler

import (
	"net/http"
	"notes/internal/errors"
	"notes/internal/models"

	"github.com/gin-gonic/gin"
)

// CreateFoodEntry — POST /food-entries
func (h *Handler) CreateFoodEntry(c *gin.Context) {
    authorID, err := h.extractAuthorID(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": errors.MsgMissingUserID})
        return
    }

    var entry models.FoodEntry
    if err := c.ShouldBindJSON(&entry); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": errors.MsgInvalidData, "details": err.Error()})
        return
    }

    entry.AuthorID = authorID

    created, err := h.service.CreateFoodEntry(c.Request.Context(), entry)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": errors.MsgDatabaseOperation, "details": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message":   errors.MsgFoodEntryCreated,
        "foodEntry": created,
    })
}

// GetFoodEntries — GET /food-entries?date=2026-03-20
func (h *Handler) GetFoodEntries(c *gin.Context) {
    authorID, err := h.extractAuthorID(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": errors.MsgMissingUserID})
        return
    }

    date := c.Query("date")

    entries, err := h.service.GetFoodEntries(c.Request.Context(), authorID, date)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": errors.MsgDatabaseOperation, "details": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message":     errors.MsgFoodEntriesFound,
        "foodEntries": entries,
        "count":       len(entries),
    })
}

// DeleteFoodEntry — DELETE /food-entries/:id
func (h *Handler) DeleteFoodEntry(c *gin.Context) {
    authorID, err := h.extractAuthorID(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": errors.MsgMissingUserID})
        return
    }

    id := c.Param("id")
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID"})
        return
    }

    if err := h.service.DeleteFoodEntry(c.Request.Context(), id, authorID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": errors.MsgDatabaseOperation, "details": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": errors.MsgFoodEntryDeleted})
}

// GetDailySummary — GET /food-entries/summary?date=2026-03-20
func (h *Handler) GetDailySummary(c *gin.Context) {
    authorID, err := h.extractAuthorID(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": errors.MsgMissingUserID})
        return
    }

    date := c.Query("date")
    if date == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Параметр date обязателен"})
        return
    }

    summary, err := h.service.GetDailySummary(c.Request.Context(), authorID, date)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": errors.MsgDatabaseOperation, "details": err.Error()})
        return
    }

    c.JSON(http.StatusOK, summary)
}

// GetWeeklySummary — GET /food-entries/summary/weekly
func (h *Handler) GetWeeklySummary(c *gin.Context) {
    authorID, err := h.extractAuthorID(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": errors.MsgMissingUserID})
        return
    }

    weekly, err := h.service.GetWeeklySummary(c.Request.Context(), authorID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": errors.MsgDatabaseOperation, "details": err.Error()})
        return
    }

    c.JSON(http.StatusOK, weekly)
}