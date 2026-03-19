package models

import "time"

// Device — регистрация FCM/APN токена
type Device struct {
    ID        string    `json:"id,omitempty" bson:"_id,omitempty"`
    AuthorID  int       `json:"authorId" bson:"author_id"`
    Token     string    `json:"token" bson:"token"`
    Platform  string    `json:"platform" bson:"platform"` // ios, android
    CreatedAt time.Time `json:"createdAt" bson:"created_at"`
}

// Notification — уведомление
type Notification struct {
    ID        string    `json:"id,omitempty" bson:"_id,omitempty"`
    AuthorID  int       `json:"authorId" bson:"author_id"`
    Title     string    `json:"title" bson:"title"`
    Body      string    `json:"body" bson:"body"`
    IsRead    bool      `json:"isRead" bson:"is_read"`
    CreatedAt time.Time `json:"createdAt" bson:"created_at"`
}