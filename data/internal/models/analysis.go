package models

import "time"

// Analysis — результаты анализов
type Analysis struct {
	ID               string    `json:"id,omitempty" bson:"_id,omitempty"`
	AuthorID         int       `json:"authorId" bson:"author_id"`
	Date             string    `json:"date" bson:"date"`
	TotalCholesterol float64   `json:"totalCholesterol" bson:"total_cholesterol"`
	LDL              float64   `json:"ldl" bson:"ldl"`
	HDL              float64   `json:"hdl" bson:"hdl"`
	Triglycerides    float64   `json:"triglycerides" bson:"triglycerides"`
	VLDL             float64   `json:"vldl" bson:"vldl"`
	CreatedAt        time.Time `json:"createdAt" bson:"created_at"`
}
