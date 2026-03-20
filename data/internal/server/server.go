package server

import (
	"fmt"

	"data/internal/config"
	"data/internal/handler"
	"data/internal/routes"
	"data/internal/service"

	"github.com/gin-gonic/gin"
)

type Server struct {
	cfg    *config.Config // Конфигурация сервера
	router *gin.Engine    // Новое поле для маршрутизатора
}

func NewServer(cfg *config.Config) (*Server, error) {
	// Проверяем, что конфигурация не пустая
	if cfg == nil {
		return nil, fmt.Errorf("конфигурация сервера не может быть nil")
	}
	service, err := service.NewService(cfg)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать сервис: %w", err)
	}
	// Создаем новый экземпляр обработчика
	handler := handler.NewHandler(cfg, service)
	// Проверяем, что обработчик успешно создан
	if handler == nil {
		return nil, fmt.Errorf("не удалось создать обработчик сервера")
	}
	fmt.Println("Обработчик сервера успешно создан")
	// Создаем новый экземпляр маршрутизатора
	router := routes.SetupRouter(handler)
	// Создаем новый экземпляр сервера
	return &Server{
		router: router,
		cfg:    cfg,
	}, nil
}

// Start - запуск сервера
func (s *Server) Start() error {
	fmt.Printf("Сервер запускается на %s:%s\n", s.cfg.Host, s.cfg.Port)
	return nil
}

// Stop - остановка сервера
func (s *Server) Stop() error {
	fmt.Println("Сервер остановлен")
	return nil
}

// Serve - основной метод сервера
func (s *Server) Serve() error {
	if err := s.Start(); err != nil {
		return err
	}

	// Запускаем сервер
	address := fmt.Sprintf("%s:%s", s.cfg.Host, s.cfg.Port)
	fmt.Printf("Сервер готов к обработке запросов на %s...\n", address)
	return s.router.Run(address)
}
