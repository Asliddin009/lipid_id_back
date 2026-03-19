package handler

import (
    "net/http"
    "notes/internal/errors"

    "github.com/gin-gonic/gin"
)

// GetLipidTrend — GET /charts/lipid-trend?period=7d
func (h *Handler) GetLipidTrend(c *gin.Context) {
    authorID, err := h.extractAuthorID(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": errors.MsgMissingUserID})
        return
    }

    period := c.DefaultQuery("period", "7d")

    trend, err := h.service.GetLipidTrend(c.Request.Context(), authorID, period)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": errors.MsgDatabaseOperation, "details": err.Error()})
        return
    }

    c.JSON(http.StatusOK, trend)
}

// GetNutritionTrend — GET /charts/nutrition-trend?period=7d
func (h *Handler) GetNutritionTrend(c *gin.Context) {
    authorID, err := h.extractAuthorID(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": errors.MsgMissingUserID})
        return
    }

    period := c.DefaultQuery("period", "7d")

    trend, err := h.service.GetNutritionTrend(c.Request.Context(), authorID, period)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": errors.MsgDatabaseOperation, "details": err.Error()})
        return
    }

    c.JSON(http.StatusOK, trend)
}

// GetWeightTrend — GET /charts/weight-trend?period=30d
func (h *Handler) GetWeightTrend(c *gin.Context) {
    authorID, err := h.extractAuthorID(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": errors.MsgMissingUserID})
        return
    }

    period := c.DefaultQuery("period", "30d")

    trend, err := h.service.GetWeightTrend(c.Request.Context(), authorID, period)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": errors.MsgDatabaseOperation, "details": err.Error()})
        return
    }

    c.JSON(http.StatusOK, trend)
}