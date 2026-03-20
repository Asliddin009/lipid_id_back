package handler

import (
	"data/internal/errors"
	"data/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateWeightEntry — POST /weight
func (h *Handler) CreateWeightEntry(c *gin.Context) {
	authorID, err := h.extractAuthorID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": errors.MsgMissingUserID})
		return
	}

	var entry models.WeightEntry
	if err := c.ShouldBindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.MsgInvalidData, "details": err.Error()})
		return
	}

	entry.AuthorID = authorID

	created, err := h.service.CreateWeightEntry(c.Request.Context(), entry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.MsgDatabaseOperation, "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     errors.MsgWeightEntryCreated,
		"weightEntry": created,
	})
}

// GetWeightEntries — GET /weight
func (h *Handler) GetWeightEntries(c *gin.Context) {
	authorID, err := h.extractAuthorID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": errors.MsgMissingUserID})
		return
	}

	entries, err := h.service.GetWeightEntries(c.Request.Context(), authorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.MsgDatabaseOperation, "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       errors.MsgWeightEntriesFound,
		"weightEntries": entries,
		"count":         len(entries),
	})
}

// DeleteWeightEntry — DELETE /weight/:id
func (h *Handler) DeleteWeightEntry(c *gin.Context) {
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

	if err := h.service.DeleteWeightEntry(c.Request.Context(), id, authorID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.MsgDatabaseOperation, "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": errors.MsgWeightEntryDeleted})
}
