package handler

import (
    "net/http"
    "notes/internal/errors"
    "notes/internal/models"

    "github.com/gin-gonic/gin"
)

// CreateAnalysis — POST /analyses
func (h *Handler) CreateAnalysis(c *gin.Context) {
    authorID, err := h.extractAuthorID(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": errors.MsgMissingUserID})
        return
    }

    var analysis models.Analysis
    if err := c.ShouldBindJSON(&analysis); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": errors.MsgInvalidData, "details": err.Error()})
        return
    }

    analysis.AuthorID = authorID

    created, err := h.service.CreateAnalysis(c.Request.Context(), analysis)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": errors.MsgDatabaseOperation, "details": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message":  errors.MsgAnalysisCreated,
        "analysis": created,
    })
}

// GetAnalyses — GET /analyses
func (h *Handler) GetAnalyses(c *gin.Context) {
    authorID, err := h.extractAuthorID(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": errors.MsgMissingUserID})
        return
    }

    analyses, err := h.service.GetAnalyses(c.Request.Context(), authorID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": errors.MsgDatabaseOperation, "details": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message":  errors.MsgAnalysesFound,
        "analyses": analyses,
        "count":    len(analyses),
    })
}

// GetAnalysisByID — GET /analyses/:id
func (h *Handler) GetAnalysisByID(c *gin.Context) {
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

    analysis, err := h.service.GetAnalysisByID(c.Request.Context(), id, authorID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": errors.MsgAnalysisFound, "details": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message":  errors.MsgAnalysisFound,
        "analysis": analysis,
    })
}

// DeleteAnalysis — DELETE /analyses/:id
func (h *Handler) DeleteAnalysis(c *gin.Context) {
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

    if err := h.service.DeleteAnalysis(c.Request.Context(), id, authorID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": errors.MsgDatabaseOperation, "details": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": errors.MsgAnalysisDeleted})
}