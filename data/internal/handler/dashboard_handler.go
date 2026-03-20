package handler

import (
	"data/internal/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetDashboard — GET /home/dashboard
func (h *Handler) GetDashboard(c *gin.Context) {
	authorID, err := h.extractAuthorID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": errors.MsgMissingUserID})
		return
	}

	dashboard, err := h.service.GetDashboard(c.Request.Context(), authorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.MsgDatabaseOperation, "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dashboard)
}
