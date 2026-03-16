package server

import (
	"auth/internal/config"
	"auth/internal/handler"
	"auth/internal/routes"
	"auth/internal/service"
	"fmt"

	"github.com/gin-gonic/gin"
)

type Server struct {
	config *config.Config
	router *gin.Engine
}

func NewServer(cfg *config.Config) (*Server, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}
	service, err := service.NewService(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create service: %w", err)
	}
	handler := handler.NewHandler(service, cfg)
	if handler == nil {
		return nil, fmt.Errorf("handler is nil")
	}
	router := routes.SetupRoutes(handler)
	if router == nil {
		return nil, fmt.Errorf("router is nil")
	}
	return &Server{
		config: cfg,
		router: router,
	}, nil
}

func (s *Server) Run() error {
	// Запускаем сервер
	address := fmt.Sprintf("%s:%s", s.config.Host, s.config.Port)
	fmt.Printf("Сервер готов к обработке запросов на %s...\n", address)
	return s.router.Run(address)
}

func (s *Server) Stop() error {
	fmt.Println("Сервер остановлен")
	return nil
}
