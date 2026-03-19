package service

import (
	"context"
	"fmt"
	"math"
	"notes/internal/caching"
	"notes/internal/config"
	"notes/internal/database"
	"notes/internal/errors"
	"notes/internal/models"
	"sort"
	"time"

	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoService struct {
	cfg   *config.Config
	db    *mongo.Client
	cache *redis.Client

	notesCol         *mongo.Collection
	foodEntriesCol   *mongo.Collection
	analysesCol      *mongo.Collection
	weightCol        *mongo.Collection
	devicesCol       *mongo.Collection
	notificationsCol *mongo.Collection
}

func NewService(cfg *config.Config) (Service, error) {
	db, err := database.NewDatabase(cfg)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrServiceCreation, err)
	}

	cache, err := caching.NewCaching(cfg)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrServiceCreation, err)
	}

	database := db.Database(cfg.DB_NAME)

	return &MongoService{
		cfg:              cfg,
		db:               db,
		cache:            cache,
		notesCol:         database.Collection(cfg.DB_COLLECTION),
		foodEntriesCol:   database.Collection("food_entries"),
		analysesCol:      database.Collection("analyses"),
		weightCol:        database.Collection("weight"),
		devicesCol:       database.Collection("devices"),
		notificationsCol: database.Collection("notifications"),
	}, nil
}

func (s *MongoService) Close() error {
	if s.cache != nil {
		if err := s.cache.Close(); err != nil {
			return fmt.Errorf("%w: %v", errors.ErrCacheClose, err)
		}
	}
	return database.CloseDB(s.db, s.cfg)
}

// ==================== Notes ====================

func (s *MongoService) Create(ctx context.Context, note models.Note) (*models.Note, error) {
	result, err := s.notesCol.InsertOne(ctx, note)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrNoteCreation, err)
	}
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		note.ID = oid.Hex()
	}
	return &note, nil
}

func (s *MongoService) GetByID(ctx context.Context, id string) (*models.Note, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrInvalidNoteID, err)
	}
	var note models.Note
	err = s.notesCol.FindOne(ctx, bson.M{"_id": objectID}).Decode(&note)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrNoteNotFound, err)
	}
	note.ID = id
	return &note, nil
}

func (s *MongoService) GetAll(ctx context.Context, authorID int) ([]models.Note, error) {
	cursor, err := s.notesCol.Find(ctx, bson.M{"author_id": authorID})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrDatabaseOperation, err)
	}
	defer cursor.Close(ctx)

	var notes []models.Note
	for cursor.Next(ctx) {
		var note models.Note
		var raw bson.M
		if err := cursor.Decode(&raw); err != nil {
			return nil, fmt.Errorf("%w: %v", errors.ErrDecodeNote, err)
		}
		if oid, ok := raw["_id"].(primitive.ObjectID); ok {
			note.ID = oid.Hex()
		}
		if v, ok := raw["name"].(string); ok {
			note.Name = v
		}
		if v, ok := raw["content"].(string); ok {
			note.Content = v
		}
		if v, ok := raw["author_id"].(int32); ok {
			note.AuthorID = int(v)
		} else if v, ok := raw["author_id"].(int64); ok {
			note.AuthorID = int(v)
		}
		notes = append(notes, note)
	}
	if notes == nil {
		notes = []models.Note{}
	}
	return notes, nil
}

func (s *MongoService) Update(ctx context.Context, note models.Note) (*models.Note, error) {
	objectID, err := primitive.ObjectIDFromHex(note.ID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrInvalidNoteID, err)
	}
	update := bson.M{
		"$set": bson.M{
			"name":      note.Name,
			"content":   note.Content,
			"author_id": note.AuthorID,
		},
	}
	_, err = s.notesCol.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrNoteUpdate, err)
	}
	return &note, nil
}

func (s *MongoService) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("%w: %v", errors.ErrInvalidNoteID, err)
	}
	_, err = s.notesCol.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("%w: %v", errors.ErrNoteDeletion, err)
	}
	return nil
}

// ==================== Food Entries ====================

func (s *MongoService) CreateFoodEntry(ctx context.Context, entry models.FoodEntry) (*models.FoodEntry, error) {
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = time.Now().UTC()
	}
	result, err := s.foodEntriesCol.InsertOne(ctx, entry)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания записи о еде: %v", err)
	}
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		entry.ID = oid.Hex()
	}
	return &entry, nil
}

func (s *MongoService) GetFoodEntries(ctx context.Context, authorID int, date string) ([]models.FoodEntry, error) {
	filter := bson.M{"author_id": authorID}

	if date != "" {
		// Парсим дату и фильтруем по дню
		t, err := time.Parse("2006-01-02", date)
		if err == nil {
			startOfDay := t.UTC()
			endOfDay := t.Add(24 * time.Hour).UTC()
			filter["created_at"] = bson.M{
				"$gte": startOfDay,
				"$lt":  endOfDay,
			}
		}
	}

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := s.foodEntriesCol.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения записей о еде: %v", err)
	}
	defer cursor.Close(ctx)

	var entries []models.FoodEntry
	if err := cursor.All(ctx, &entries); err != nil {
		return nil, fmt.Errorf("ошибка декодирования записей о еде: %v", err)
	}

	// Заполняем ID из _id
	for i := range entries {
		if entries[i].ID == "" {
			// Попробуем получить через raw decode
		}
	}

	if entries == nil {
		entries = []models.FoodEntry{}
	}
	return entries, nil
}

func (s *MongoService) DeleteFoodEntry(ctx context.Context, id string, authorID int) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("некорректный ID записи о еде: %v", err)
	}
	result, err := s.foodEntriesCol.DeleteOne(ctx, bson.M{"_id": objectID, "author_id": authorID})
	if err != nil {
		return fmt.Errorf("ошибка удаления записи о еде: %v", err)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("запись о еде не найдена или нет доступа")
	}
	return nil
}

func (s *MongoService) GetDailySummary(ctx context.Context, authorID int, date string) (*models.DailySummary, error) {
	entries, err := s.GetFoodEntries(ctx, authorID, date)
	if err != nil {
		return nil, err
	}

	summary := &models.DailySummary{
		Date: date,
		ByMealTime: map[string]int{
			"breakfast": 0,
			"lunch":     0,
			"snack":     0,
			"dinner":    0,
		},
	}

	for _, e := range entries {
		summary.TotalCalories += e.Calories
		summary.TotalProteins += e.Proteins
		summary.TotalFats += e.Fats
		summary.TotalCarbs += e.Carbs
		if _, ok := summary.ByMealTime[e.MealTime]; ok {
			summary.ByMealTime[e.MealTime] += e.Calories
		}
	}

	// Округляем
	summary.TotalProteins = math.Round(summary.TotalProteins*10) / 10
	summary.TotalFats = math.Round(summary.TotalFats*10) / 10
	summary.TotalCarbs = math.Round(summary.TotalCarbs*10) / 10

	return summary, nil
}

func (s *MongoService) GetWeeklySummary(ctx context.Context, authorID int) (*models.WeeklySummary, error) {
	now := time.Now().UTC()
	weekly := &models.WeeklySummary{
		Days: make([]models.DailySummary, 0, 7),
	}

	for i := 6; i >= 0; i-- {
		day := now.AddDate(0, 0, -i)
		dateStr := day.Format("2006-01-02")
		summary, err := s.GetDailySummary(ctx, authorID, dateStr)
		if err != nil {
			// Возвращаем пустой день при ошибке
			summary = &models.DailySummary{
				Date:       dateStr,
				ByMealTime: map[string]int{"breakfast": 0, "lunch": 0, "snack": 0, "dinner": 0},
			}
		}
		weekly.Days = append(weekly.Days, *summary)
	}

	return weekly, nil
}

// ==================== Dashboard ====================

func (s *MongoService) GetDashboard(ctx context.Context, authorID int) (*models.DashboardResponse, error) {
	// Собираем данные для дашборда

	// 1. Липидный скор — берём последний анализ
	score := 0.0
	analyses, err := s.GetAnalyses(ctx, authorID)
	if err == nil && len(analyses) > 0 {
		latest := analyses[0]
		// Простая формула скора на основе анализов (можно заменить на более сложную)
		score = s.calculateLipidScore(latest)
	}

	// 2. Chart data — последние 7 дней липидного скора
	chartData := s.buildChartData(ctx, authorID, 7)

	// 3. Recent events
	recentEvents := s.buildRecentEvents(ctx, authorID)

	return &models.DashboardResponse{
		Score:        score,
		ChartData:    chartData,
		RecentEvents: recentEvents,
	}, nil
}

func (s *MongoService) calculateLipidScore(a models.Analysis) float64 {
	// Упрощённая формула: 10 - (LDL - 2.5) - (Triglycerides - 1.0) + (HDL - 1.0)
	score := 10.0 - (a.LDL-2.5) - (a.Triglycerides-1.0) + (a.HDL-1.0)
	if score < 0 {
		score = 0
	}
	if score > 10 {
		score = 10
	}
	return math.Round(score*10) / 10
}

func (s *MongoService) buildChartData(ctx context.Context, authorID int, days int) []models.ChartPoint {
	points := make([]models.ChartPoint, 0, days)
	now := time.Now().UTC()

	// Получаем все анализы пользователя
	analyses, _ := s.GetAnalyses(ctx, authorID)

	// Создаём карту анализов по дате
	analysisMap := make(map[string]models.Analysis)
	for _, a := range analyses {
		analysisMap[a.Date] = a
	}

	lastScore := 5.0 // Начальный скор по умолчанию
	if len(analyses) > 0 {
		lastScore = s.calculateLipidScore(analyses[len(analyses)-1])
	}

	for i := days - 1; i >= 0; i-- {
		day := now.AddDate(0, 0, -i)
		dateStr := day.Format("2006-01-02")
		if a, ok := analysisMap[dateStr]; ok {
			lastScore = s.calculateLipidScore(a)
		}
		points = append(points, models.ChartPoint{
			Date:  dateStr,
			Value: lastScore,
		})
	}

	return points
}

func (s *MongoService) buildRecentEvents(ctx context.Context, authorID int) []models.RecentEvent {
	events := make([]models.RecentEvent, 0)
	now := time.Now().UTC()
	weekAgo := now.AddDate(0, 0, -7)

	// Последние записи о еде
	foodFilter := bson.M{
		"author_id":  authorID,
		"created_at": bson.M{"$gte": weekAgo},
	}
	foodOpts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(3)
	foodCursor, err := s.foodEntriesCol.Find(ctx, foodFilter, foodOpts)
	if err == nil {
		var foodEntries []models.FoodEntry
		if err := foodCursor.All(ctx, &foodEntries); err == nil {
			for _, fe := range foodEntries {
				events = append(events, models.RecentEvent{
					Type:      "meal",
					Title:     "Приём пищи",
					Subtitle:  fe.MealTime + " — " + fe.DishName,
					Trailing:  "",
					Color:     "blue",
					CreatedAt: fe.CreatedAt,
				})
			}
		}
		foodCursor.Close(ctx)
	}

	// Последние анализы
	analysisFilter := bson.M{
		"author_id":  authorID,
		"created_at": bson.M{"$gte": weekAgo},
	}
	analysisOpts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(3)
	analysisCursor, err := s.analysesCol.Find(ctx, analysisFilter, analysisOpts)
	if err == nil {
		var analysisEntries []models.Analysis
		if err := analysisCursor.All(ctx, &analysisEntries); err == nil {
			for _, a := range analysisEntries {
				events = append(events, models.RecentEvent{
					Type:      "lab",
					Title:     "Лабораторные измерения",
					Subtitle:  fmt.Sprintf("LDL: %.1f ммоль/л", a.LDL),
					Trailing:  "",
					Color:     "green",
					CreatedAt: a.CreatedAt,
				})
			}
		}
		analysisCursor.Close(ctx)
	}

	// Последние замеры веса
	weightFilter := bson.M{
		"author_id":  authorID,
		"created_at": bson.M{"$gte": weekAgo},
	}
	weightOpts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(3)
	weightCursor, err := s.weightCol.Find(ctx, weightFilter, weightOpts)
	if err == nil {
		var weightEntries []models.WeightEntry
		if err := weightCursor.All(ctx, &weightEntries); err == nil {
			for _, w := range weightEntries {
				events = append(events, models.RecentEvent{
					Type:      "anthropometry",
					Title:     "Антропометрия",
					Subtitle:  fmt.Sprintf("%.1f кг", w.Value),
					Trailing:  "",
					Color:     "gray",
					CreatedAt: w.CreatedAt,
				})
			}
		}
		weightCursor.Close(ctx)
	}

	// Сортируем по дате убывания
	sort.Slice(events, func(i, j int) bool {
		return events[i].CreatedAt.After(events[j].CreatedAt)
	})

	// Ограничиваем до 10 событий
	if len(events) > 10 {
		events = events[:10]
	}

	return events
}

// ==================== Charts ====================

func (s *MongoService) parsePeriodDays(period string) int {
	switch period {
	case "7d":
		return 7
	case "30d":
		return 30
	case "90d":
		return 90
	default:
		return 7
	}
}

func (s *MongoService) GetLipidTrend(ctx context.Context, authorID int, period string) (*models.ChartTrendResponse, error) {
	days := s.parsePeriodDays(period)
	chartData := s.buildChartData(ctx, authorID, days)
	return &models.ChartTrendResponse{
		Period: period,
		Points: chartData,
	}, nil
}

func (s *MongoService) GetNutritionTrend(ctx context.Context, authorID int, period string) (*models.ChartTrendResponse, error) {
	days := s.parsePeriodDays(period)
	now := time.Now().UTC()
	points := make([]models.ChartPoint, 0, days)

	for i := days - 1; i >= 0; i-- {
		day := now.AddDate(0, 0, -i)
		dateStr := day.Format("2006-01-02")
		summary, _ := s.GetDailySummary(ctx, authorID, dateStr)
		calories := 0.0
		if summary != nil {
			calories = float64(summary.TotalCalories)
		}
		points = append(points, models.ChartPoint{
			Date:  dateStr,
			Value: calories,
		})
	}

	return &models.ChartTrendResponse{
		Period: period,
		Points: points,
	}, nil
}

func (s *MongoService) GetWeightTrend(ctx context.Context, authorID int, period string) (*models.ChartTrendResponse, error) {
	days := s.parsePeriodDays(period)
	now := time.Now().UTC()
	startDate := now.AddDate(0, 0, -days)

	filter := bson.M{
		"author_id":  authorID,
		"created_at": bson.M{"$gte": startDate},
	}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: 1}})
	cursor, err := s.weightCol.Find(ctx, filter, opts)
	if err != nil {
		return &models.ChartTrendResponse{Period: period, Points: []models.ChartPoint{}}, nil
	}
	defer cursor.Close(ctx)

	var entries []models.WeightEntry
	if err := cursor.All(ctx, &entries); err != nil {
		return &models.ChartTrendResponse{Period: period, Points: []models.ChartPoint{}}, nil
	}

	points := make([]models.ChartPoint, 0, len(entries))
	for _, e := range entries {
		points = append(points, models.ChartPoint{
			Date:  e.Date,
			Value: e.Value,
		})
	}

	return &models.ChartTrendResponse{
		Period: period,
		Points: points,
	}, nil
}

// ==================== Analyses ====================

func (s *MongoService) CreateAnalysis(ctx context.Context, analysis models.Analysis) (*models.Analysis, error) {
	if analysis.CreatedAt.IsZero() {
		analysis.CreatedAt = time.Now().UTC()
	}
	result, err := s.analysesCol.InsertOne(ctx, analysis)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания анализа: %v", err)
	}
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		analysis.ID = oid.Hex()
	}
	return &analysis, nil
}

func (s *MongoService) GetAnalyses(ctx context.Context, authorID int) ([]models.Analysis, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := s.analysesCol.Find(ctx, bson.M{"author_id": authorID}, opts)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения анализов: %v", err)
	}
	defer cursor.Close(ctx)

	var analyses []models.Analysis
	if err := cursor.All(ctx, &analyses); err != nil {
		return nil, fmt.Errorf("ошибка декодирования анализов: %v", err)
	}
	if analyses == nil {
		analyses = []models.Analysis{}
	}
	return analyses, nil
}

func (s *MongoService) GetAnalysisByID(ctx context.Context, id string, authorID int) (*models.Analysis, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("некорректный ID анализа: %v", err)
	}
	var analysis models.Analysis
	err = s.analysesCol.FindOne(ctx, bson.M{"_id": objectID, "author_id": authorID}).Decode(&analysis)
	if err != nil {
		return nil, fmt.Errorf("анализ не найден: %v", err)
	}
	analysis.ID = id
	return &analysis, nil
}

func (s *MongoService) DeleteAnalysis(ctx context.Context, id string, authorID int) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("некорректный ID анализа: %v", err)
	}
	result, err := s.analysesCol.DeleteOne(ctx, bson.M{"_id": objectID, "author_id": authorID})
	if err != nil {
		return fmt.Errorf("ошибка удаления анализа: %v", err)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("анализ не найден или нет доступа")
	}
	return nil
}

// ==================== Weight ====================

func (s *MongoService) CreateWeightEntry(ctx context.Context, entry models.WeightEntry) (*models.WeightEntry, error) {
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = time.Now().UTC()
	}
	result, err := s.weightCol.InsertOne(ctx, entry)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания замера веса: %v", err)
	}
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		entry.ID = oid.Hex()
	}
	return &entry, nil
}

func (s *MongoService) GetWeightEntries(ctx context.Context, authorID int) ([]models.WeightEntry, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := s.weightCol.Find(ctx, bson.M{"author_id": authorID}, opts)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения замеров веса: %v", err)
	}
	defer cursor.Close(ctx)

	var entries []models.WeightEntry
	if err := cursor.All(ctx, &entries); err != nil {
		return nil, fmt.Errorf("ошибка декодирования замеров веса: %v", err)
	}
	if entries == nil {
		entries = []models.WeightEntry{}
	}
	return entries, nil
}

func (s *MongoService) DeleteWeightEntry(ctx context.Context, id string, authorID int) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("некорректный ID замера веса: %v", err)
	}
	result, err := s.weightCol.DeleteOne(ctx, bson.M{"_id": objectID, "author_id": authorID})
	if err != nil {
		return fmt.Errorf("ошибка удаления замера веса: %v", err)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("замер веса не найден или нет доступа")
	}
	return nil
}

// ==================== Devices & Notifications ====================

func (s *MongoService) RegisterDevice(ctx context.Context, device models.Device) (*models.Device, error) {
	if device.CreatedAt.IsZero() {
		device.CreatedAt = time.Now().UTC()
	}
	// Upsert: если токен уже существует, обновляем
	filter := bson.M{"author_id": device.AuthorID, "token": device.Token}
	update := bson.M{"$set": device}
	opts := options.Update().SetUpsert(true)
	result, err := s.devicesCol.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return nil, fmt.Errorf("ошибка регистрации устройства: %v", err)
	}
	if result.UpsertedID != nil {
		if oid, ok := result.UpsertedID.(primitive.ObjectID); ok {
			device.ID = oid.Hex()
		}
	}
	return &device, nil
}

func (s *MongoService) GetNotifications(ctx context.Context, authorID int) ([]models.Notification, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := s.notificationsCol.Find(ctx, bson.M{"author_id": authorID}, opts)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения уведомлений: %v", err)
	}
	defer cursor.Close(ctx)

	var notifications []models.Notification
	if err := cursor.All(ctx, &notifications); err != nil {
		return nil, fmt.Errorf("ошибка декодирования уведомлений: %v", err)
	}
	if notifications == nil {
		notifications = []models.Notification{}
	}
	return notifications, nil
}

func (s *MongoService) MarkNotificationRead(ctx context.Context, id string, authorID int) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("некорректный ID уведомления: %v", err)
	}
	result, err := s.notificationsCol.UpdateOne(
		ctx,
		bson.M{"_id": objectID, "author_id": authorID},
		bson.M{"$set": bson.M{"is_read": true}},
	)
	if err != nil {
		return fmt.Errorf("ошибка обновления уведомления: %v", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("уведомление не найдено или нет доступа")
	}
	return nil
}
