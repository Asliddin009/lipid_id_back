package models

import "time"

// FoodEntry — запись о приёме пищи
type FoodEntry struct {
    ID        string    `json:"id,omitempty" bson:"_id,omitempty"`
    AuthorID  int       `json:"authorId" bson:"author_id"`
    DishName  string    `json:"dishName" bson:"dish_name"`
    MealTime  string    `json:"mealTime" bson:"meal_time"` // breakfast, lunch, snack, dinner
    Calories  int       `json:"calories" bson:"calories"`
    Proteins  float64   `json:"proteins" bson:"proteins"`
    Fats      float64   `json:"fats" bson:"fats"`
    Carbs     float64   `json:"carbs" bson:"carbs"`
    ImageUrl  *string   `json:"imageUrl" bson:"image_url,omitempty"`
    CreatedAt time.Time `json:"createdAt" bson:"created_at"`
}

// DailySummary — дневная статистика КБЖУ
type DailySummary struct {
    Date          string         `json:"date"`
    TotalCalories int            `json:"totalCalories"`
    TotalProteins float64        `json:"totalProteins"`
    TotalFats     float64        `json:"totalFats"`
    TotalCarbs    float64        `json:"totalCarbs"`
    ByMealTime    map[string]int `json:"byMealTime"`
}

// WeeklySummary — недельная статистика
type WeeklySummary struct {
    Days []DailySummary `json:"days"`
}