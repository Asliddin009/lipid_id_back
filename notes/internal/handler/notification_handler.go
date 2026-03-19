package handler

import (
    "net/http"
    "notes/internal/errors"
    "notes/internal/models"

    "github.com/gin-gonic/gin"
)

// RegisterDevice — POST /devices
func (h *Handler) RegisterDevice(c *gin.Context) {
    authorID, err := h.extractAuthorID(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": errors.MsgMissingUserID})
        return
    }

    var device models.Device
    if err := c.ShouldBindJSON(&device); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": errors.MsgInvalidData, "details": err.Error()})
        return
    }

    device.AuthorID = authorID

    created, err := h.service.RegisterDevice(c.Request.Context(), device)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": errors.MsgDatabaseOperation, "details": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": errors.MsgDeviceRegistered,
        "device":  created,
    })
}

// GetNotifications — GET /notifications
func (h *Handler) GetNotifications(c *gin.Context) {
    authorID, err := h.extractAuthorID(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": errors.MsgMissingUserID})
        return
    }

    notifications, err := h.service.GetNotifications(c.Request.Context(), authorID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": errors.MsgDatabaseOperation, "details": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message":       errors.MsgNotificationsFound,
        "notifications": notifications,
        "count":         len(notifications),
    })
}

// MarkNotificationRead — PATCH /notifications/:id/read
func (h *Handler) MarkNotificationRead(c *gin.Context) {
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

    if err := h.service.MarkNotificationRead(c.Request.Context(), id, authorID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": errors.MsgDatabaseOperation, "details": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": errors.MsgNotificationMarkedRead})
}