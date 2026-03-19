package service

import (
	"context"
	"notes/internal/models"
)

type Service interface {
	Close() error                                                       // Закрывает соединение с базой данных
	Create(ctx context.Context, note models.Note) (*models.Note, error) // Создает новую запись в базе данных
	GetByID(ctx context.Context, id string) (*models.Note, error)       // Получает запись по идентификатору
	GetAll(ctx context.Context, authorId int) ([]models.Note, error)    // Получает все заметки из базы данных
	Update(ctx context.Context, note models.Note) (*models.Note, error) // Обновляет существующую запись
	Delete(ctx context.Context, id string) error                        // Удаляет запись по идентификатору
}
