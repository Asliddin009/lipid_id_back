package service

import (
	"context"
	"notes/internal/models"
)

type Service interface {
	Close() error // Закрывает соединение с базой данных

	// === Notes ===
	Create(ctx context.Context, note models.Note) (*models.Note, error)
	GetByID(ctx context.Context, id string) (*models.Note, error)
	GetAll(ctx context.Context, authorId int) ([]models.Note, error)
	Update(ctx context.Context, note models.Note) (*models.Note, error)
	Delete(ctx context.Context, id string) error

	// === Food Entries ===
	CreateFoodEntry(ctx context.Context, entry models.FoodEntry) (*models.FoodEntry, error)
	GetFoodEntries(ctx context.Context, authorID int, date string) ([]models.FoodEntry, error)
	DeleteFoodEntry(ctx context.Context, id string, authorID int) error
	GetDailySummary(ctx context.Context, authorID int, date string) (*models.DailySummary, error)
	GetWeeklySummary(ctx context.Context, authorID int) (*models.WeeklySummary, error)

	// === Dashboard ===
	GetDashboard(ctx context.Context, authorID int) (*models.DashboardResponse, error)

	// === Charts ===
	GetLipidTrend(ctx context.Context, authorID int, period string) (*models.ChartTrendResponse, error)
	GetNutritionTrend(ctx context.Context, authorID int, period string) (*models.ChartTrendResponse, error)
	GetWeightTrend(ctx context.Context, authorID int, period string) (*models.ChartTrendResponse, error)

	// === Analyses ===
	CreateAnalysis(ctx context.Context, analysis models.Analysis) (*models.Analysis, error)
	GetAnalyses(ctx context.Context, authorID int) ([]models.Analysis, error)
	GetAnalysisByID(ctx context.Context, id string, authorID int) (*models.Analysis, error)
	DeleteAnalysis(ctx context.Context, id string, authorID int) error

	// === Weight ===
	CreateWeightEntry(ctx context.Context, entry models.WeightEntry) (*models.WeightEntry, error)
	GetWeightEntries(ctx context.Context, authorID int) ([]models.WeightEntry, error)
	DeleteWeightEntry(ctx context.Context, id string, authorID int) error

	// === Devices & Notifications ===
	RegisterDevice(ctx context.Context, device models.Device) (*models.Device, error)
	GetNotifications(ctx context.Context, authorID int) ([]models.Notification, error)
	MarkNotificationRead(ctx context.Context, id string, authorID int) error
}
