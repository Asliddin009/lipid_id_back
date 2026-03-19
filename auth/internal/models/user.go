package models

import (
	"math"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                    int     `json:"id" gorm:"primaryKey"`
	Username              string  `json:"username" gorm:"unique;not null"`
	Password              string  `json:"password,omitempty" gorm:"not null"`
	Gender                string  `json:"gender" gorm:"default:''"`
	Age                   int     `json:"age" gorm:"default:0"`
	Height                float64 `json:"height" gorm:"default:0"`
	Weight                float64 `json:"weight" gorm:"default:0"`
	BMI                   float64 `json:"bmi" gorm:"default:0"`
	Goal                  string  `json:"goal" gorm:"default:''"`
	ActivityLevel         string  `json:"activityLevel" gorm:"column:activity_level;default:''"`
	Language              string  `json:"language" gorm:"default:'ru'"`
	NotificationsEnabled  bool    `json:"notificationsEnabled" gorm:"column:notifications_enabled;default:true"`
}

const bcryptCost = 12

// Метод для хеширования пароля
func (u *User) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// Метод для проверки пароля
func (u *User) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// CalculateBMI вычисляет BMI на основе веса (кг) и роста (см)
func (u *User) CalculateBMI() {
	if u.Height > 0 && u.Weight > 0 {
		heightM := u.Height / 100.0
		u.BMI = math.Round(u.Weight/(heightM*heightM)*10) / 10
	}
}
