package models

import "time"

// DashboardResponse — ответ главной страницы
type DashboardResponse struct {
	Score        float64       `json:"score"`
	ChartData    []ChartPoint  `json:"chartData"`
	RecentEvents []RecentEvent `json:"recentEvents"`
}

// ChartPoint — точка на графике
type ChartPoint struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
}

// RecentEvent — недавнее событие
type RecentEvent struct {
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Subtitle  string    `json:"subtitle"`
	Trailing  string    `json:"trailing"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"createdAt"`
}

// ChartTrendResponse — ответ для графиков трендов
type ChartTrendResponse struct {
	Period string       `json:"period"`
	Points []ChartPoint `json:"points"`
}
