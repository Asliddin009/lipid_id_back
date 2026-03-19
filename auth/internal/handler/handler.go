package handler

import (
	"auth/internal/config"
	"auth/internal/errors"
	"auth/internal/models"
	"auth/internal/service"
	"context"
	"time"

	jwtmanager "jwt_manager"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	cfg        *config.Config
	jwtManager *jwtmanager.JWTManager
	service    service.Service
}

func NewHandler(service service.Service, config *config.Config) *Handler {
	// Создаем JWT менеджер
	jwtConfig := jwtmanager.JWTConfig{
		SecretKey:              config.JWTSecretKey,
		AccessTokenExpiration:  config.AccessTokenExpiration,
		RefreshTokenExpiration: config.RefreshTokenExpiration,
	}
	jwtManager := jwtmanager.NewJWTManager(jwtConfig)

	return &Handler{
		cfg:        config,
		jwtManager: jwtManager,
		service:    service,
	}

}

// RegisterUser обрабатывает запрос на регистрацию нового пользователя
func (h *Handler) RegisterUser(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{
			"error":   errors.MsgInvalidData,
			"details": err.Error(),
		})
		return
	}

	if user.Username == "" || user.Password == "" {
		c.JSON(400, gin.H{
			"error": errors.MsgInvalidUserData,
		})
		return
	}

	// Пересчитываем BMI если переданы рост и вес
	user.CalculateBMI()

	// Создаем пользователя в базе данных с таймаутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.cfg.DBTimeout)*time.Second)
	defer cancel()

	createdUser, err := h.service.Create(ctx, &user)
	if err != nil {
		c.JSON(500, gin.H{
			"error":   errors.MsgUserCreation,
			"details": err.Error(),
		})
		return
	}

	createdUser.Password = ""

	c.JSON(201, gin.H{
		"message": errors.MsgUserRegistered,
		"user":    createdUser.ToProfileResponse(),
	})
}

func (h *Handler) LoginUser(c *gin.Context) {
	var loginRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(400, gin.H{
			"error":   errors.MsgInvalidData,
			"details": err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.cfg.DBTimeout)*time.Second)

	defer cancel()

	user, err := h.service.Authenticate(ctx, loginRequest.Username, loginRequest.Password)
	if err != nil {
		c.JSON(401, gin.H{
			"error": errors.MsgInvalidCredentials,
		})
		return
	}

	user.Password = ""

	accessToken, refreshToken, err := h.jwtManager.GenerateTokens(user.ID)
	if err != nil {
		c.JSON(500, gin.H{
			"error":   errors.MsgTokenGeneration,
			"details": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message":       errors.MsgLoginSuccess,
		"user":          user,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// GetUserInfo обрабатывает запрос на получение информации о пользователе
func (h *Handler) GetUserInfo(c *gin.Context) {

	userID, err := h.GetCurrentUserID(c)
	if err != nil {
		c.JSON(401, gin.H{
			"error": errors.MsgAuthRequired,
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.cfg.DBTimeout)*time.Second)
	defer cancel()

	user, err := h.service.Read(ctx, userID)
	if err != nil {
		c.JSON(404, gin.H{
			"error": errors.MsgUserNotFound,
		})
		return
	}

	c.JSON(200, user.ToProfileResponse())
}

// GetCurrentUserID получает ID текущего пользователя из контекста (для совместимости)
func (h *Handler) GetCurrentUserID(c *gin.Context) (int, error) {
	return jwtmanager.GetCurrentUserID(c)
}

func (h *Handler) UpdateUser(c *gin.Context) {
	userID, err := h.GetCurrentUserID(c)
	if err != nil {
		c.JSON(401, gin.H{
			"error": errors.MsgAuthRequired,
		})
		return
	}

	var updateReq models.UpdateProfileRequest

	// Парсим JSON из тела запроса
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(400, gin.H{
			"error":   errors.MsgInvalidData,
			"details": err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.cfg.DBTimeout)*time.Second)
	defer cancel()

	// Получаем текущего пользователя
	user, err := h.service.Read(ctx, userID)
	if err != nil {
		c.JSON(404, gin.H{
			"error": errors.MsgUserNotFound,
		})
		return
	}

	// Применяем частичное обновление
	user.ApplyUpdate(&updateReq)

	// Если передан новый пароль — хешируем
	if updateReq.Password != nil && *updateReq.Password != "" {
		hashedPassword, err := user.HashPassword(*updateReq.Password)
		if err != nil {
			c.JSON(500, gin.H{
				"error": "Ошибка хеширования пароля",
			})
			return
		}
		user.Password = hashedPassword
	}

	updateCtx, updateCancel := context.WithTimeout(c.Request.Context(), time.Duration(h.cfg.DBTimeout)*time.Second)
	defer updateCancel()

	err = h.service.Update(updateCtx, user)
	if err != nil {
		c.JSON(500, gin.H{
			"error":   errors.MsgDatabaseOperation,
			"details": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": errors.MsgUserUpdated,
		"user":    user.ToProfileResponse(),
	})
}

// DeleteUser обрабатывает запрос на удаление пользователя
func (h *Handler) DeleteUser(c *gin.Context) {
	userID, err := h.GetCurrentUserID(c)
	if err != nil {
		c.JSON(401, gin.H{
			"error": errors.MsgAuthRequired,
		})
		return
	}

	checkCtx, checkCancel := context.WithTimeout(c.Request.Context(), time.Duration(h.cfg.DBTimeout)*time.Second)
	defer checkCancel()

	_, err = h.service.Read(checkCtx, userID)
	if err != nil {
		c.JSON(404, gin.H{
			"error": errors.MsgUserNotFound,
		})
		return
	}

	deleteCtx, deleteCancel := context.WithTimeout(c.Request.Context(), time.Duration(h.cfg.DBTimeout)*time.Second)
	defer deleteCancel()

	err = h.service.Delete(deleteCtx, userID)
	if err != nil {
		c.JSON(500, gin.H{
			"error":   errors.MsgDatabaseOperation,
			"details": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": errors.MsgUserDeleted,
	})
}

// RefreshToken обрабатывает запрос на обновление access токена с помощью refresh токена
func (h *Handler) RefreshToken(c *gin.Context) {
	var refreshRequest struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	// Парсим JSON из тела запроса
	if err := c.ShouldBindJSON(&refreshRequest); err != nil {
		c.JSON(400, gin.H{
			"error":   errors.MsgInvalidData,
			"details": err.Error(),
		})
		return
	}

	// Валидируем refresh токен
	userID, err := h.jwtManager.ValidateRefreshToken(refreshRequest.RefreshToken)
	if err != nil {
		c.JSON(401, gin.H{
			"error": errors.MsgRefreshToken,
		})
		return
	}

	// Проверяем, что пользователь существует
	user, err := h.service.Read(c.Request.Context(), userID)
	if err != nil {
		c.JSON(404, gin.H{
			"error": errors.MsgUserNotFound,
		})
		return
	}

	// Генерируем новые токены
	accessToken, refreshToken, err := h.jwtManager.GenerateTokens(user.ID)
	if err != nil {
		c.JSON(500, gin.H{
			"error":   errors.MsgTokenGeneration,
			"details": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message":       errors.MsgTokensRefreshed,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// ExtractTokenFromHeader извлекает JWT токен из HTTP заголовка Authorization (публичный метод)
func (h *Handler) ExtractTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.ErrMissingAuthHeader
	}

	// Проверяем на формат Bearer
	const bearerPrefix = "Bearer "
	if len(authHeader) <= len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		return "", errors.ErrInvalidAuthFormat
	}

	// Возвращаем токен без префикса "Bearer "
	return authHeader[len(bearerPrefix):], nil
}

// ValidateAccessToken валидирует access токен (публичный метод)
func (h *Handler) ValidateAccessToken(tokenString string) (int, error) {
	return h.jwtManager.ValidateAccessToken(tokenString)
}

// RequireAuth возвращает middleware для проверки JWT токена
func (h *Handler) RequireAuth() gin.HandlerFunc {
	return h.jwtManager.JWTInterceptor()
}
