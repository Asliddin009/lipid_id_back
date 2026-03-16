package service

import (
	"auth/internal/models"
	"context"
)

type Service interface {
	// Создает нового пользователя
	Create(ctx context.Context, user *models.User) (*models.User, error)
	// Находит пользователя по ID
	Read(ctx context.Context, id int) (*models.User, error)
	// Обновляет данные пользователя
	Update(ctx context.Context, user *models.User) error
	// Удаляет пользователя по ID
	Delete(ctx context.Context, id int) error
	// Аутентифицирует пользователя по имени пользователя и паролю
	Authenticate(ctx context.Context, username, password string) (*models.User, error)
	// Находит пользователя по имени пользователя
	ReadByUsername(ctx context.Context, username string) (*models.User, error)
	// Close закрывает соединение с базой данных
	Close() error
}
