package handler

// import (
// 	"auth/internal/config"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/gin-gonic/gin"
// 	"github.com/stretchr/testify/assert"
// )

// // setupRouter создает новый gin-роутер для тестирования и регистрирует обработчики.
// func setupRouter() (*gin.Engine, *Handler) {
// 	// Используем тестовый режим Gin
// 	gin.SetMode(gin.TestMode)

// 	// Создаем моковую конфигурацию
// 	mockConfig := &config.Config{}
// 	handler := NewHandler(mockConfig)

// 	router := gin.Default()
// 	authRoutes := router.Group("/auth")
// 	{
// 		authRoutes.POST("/register", handler.RegisterUser)
// 		authRoutes.POST("/login", handler.LoginUser)
// 		authRoutes.GET("/user", handler.GetUserInfo)
// 		authRoutes.PUT("/user", handler.UpdateUser)
// 		authRoutes.DELETE("/user", handler.DeleteUser)
// 	}

// 	return router, handler
// }

// func TestRegisterUser(t *testing.T) {
// 	router, _ := setupRouter()

// 	w := httptest.NewRecorder()
// 	req, _ := http.NewRequest(http.MethodPost, "/auth/register", nil)
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusCreated, w.Code)

// 	var response map[string]string
// 	err := json.Unmarshal(w.Body.Bytes(), &response)
// 	assert.NoError(t, err)
// 	assert.Equal(t, "User registered successfully", response["message"])
// }

// func TestLoginUser(t *testing.T) {
// 	router, _ := setupRouter()

// 	w := httptest.NewRecorder()
// 	req, _ := http.NewRequest(http.MethodPost, "/auth/login", nil)
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusOK, w.Code)

// 	var response map[string]string
// 	err := json.Unmarshal(w.Body.Bytes(), &response)
// 	assert.NoError(t, err)
// 	assert.Equal(t, "Пользователь успешно авторизован", response["message"])
// 	assert.Equal(t, "example_access_token", response["access_token"])
// 	assert.Equal(t, "example_refresh_token", response["refresh_token"])
// }

// func TestGetUserInfo(t *testing.T) {
// 	router, _ := setupRouter()

// 	w := httptest.NewRecorder()
// 	req, _ := http.NewRequest(http.MethodGet, "/auth/user", nil)
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusOK, w.Code)

// 	var response map[string]string
// 	err := json.Unmarshal(w.Body.Bytes(), &response)
// 	assert.NoError(t, err)
// 	assert.Equal(t, "User retrieved successfully", response["message"])
// }

// func TestUpdateUser(t *testing.T) {
// 	router, _ := setupRouter()

// 	w := httptest.NewRecorder()
// 	req, _ := http.NewRequest(http.MethodPut, "/auth/user", nil)
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusOK, w.Code)

// 	var response map[string]string
// 	err := json.Unmarshal(w.Body.Bytes(), &response)
// 	assert.NoError(t, err)
// 	assert.Equal(t, "User updated successfully", response["message"])
// }

// func TestDeleteUser(t *testing.T) {
// 	router, _ := setupRouter()

// 	w := httptest.NewRecorder()
// 	req, _ := http.NewRequest(http.MethodDelete, "/auth/user", nil)
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusOK, w.Code)

// 	var response map[string]string
// 	err := json.Unmarshal(w.Body.Bytes(), &response)
// 	assert.NoError(t, err)
// 	assert.Equal(t, "User deleted successfully", response["message"])
// }
