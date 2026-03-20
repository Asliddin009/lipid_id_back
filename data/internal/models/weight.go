package models

import "time"

// WeightEntry — замер веса
type WeightEntry struct {
	ID        string    `json:"id,omitempty" bson:"_id,omitempty"`
	AuthorID  int       `json:"authorId" bson:"author_id"`
	Value     float64   `json:"value" bson:"value"`
	Date      string    `json:"date" bson:"date"`
	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
}
