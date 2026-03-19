package models

// ProfileResponse — ответ GET /profile (без пароля)
type ProfileResponse struct {
    ID                   int     `json:"id"`
    Username             string  `json:"username"`
    Gender               string  `json:"gender"`
    Age                  int     `json:"age"`
    Height               float64 `json:"height"`
    Weight               float64 `json:"weight"`
    BMI                  float64 `json:"bmi"`
    Goal                 string  `json:"goal"`
    ActivityLevel        string  `json:"activityLevel"`
    Language             string  `json:"language"`
    NotificationsEnabled bool    `json:"notificationsEnabled"`
}

// UpdateProfileRequest — тело запроса для обновления профиля
type UpdateProfileRequest struct {
    Username             *string  `json:"username,omitempty"`
    Password             *string  `json:"password,omitempty"`
    Gender               *string  `json:"gender,omitempty"`
    Age                  *int     `json:"age,omitempty"`
    Height               *float64 `json:"height,omitempty"`
    Weight               *float64 `json:"weight,omitempty"`
    Goal                 *string  `json:"goal,omitempty"`
    ActivityLevel        *string  `json:"activityLevel,omitempty"`
    Language             *string  `json:"language,omitempty"`
    NotificationsEnabled *bool    `json:"notificationsEnabled,omitempty"`
}

// ToProfileResponse конвертирует User в ProfileResponse
func (u *User) ToProfileResponse() *ProfileResponse {
    return &ProfileResponse{
        ID:                   u.ID,
        Username:             u.Username,
        Gender:               u.Gender,
        Age:                  u.Age,
        Height:               u.Height,
        Weight:               u.Weight,
        BMI:                  u.BMI,
        Goal:                 u.Goal,
        ActivityLevel:        u.ActivityLevel,
        Language:             u.Language,
        NotificationsEnabled: u.NotificationsEnabled,
    }
}

// ApplyUpdate применяет частичное обновление к пользователю
func (u *User) ApplyUpdate(req *UpdateProfileRequest) {
    if req.Username != nil {
        u.Username = *req.Username
    }
    if req.Gender != nil {
        u.Gender = *req.Gender
    }
    if req.Age != nil {
        u.Age = *req.Age
    }
    if req.Height != nil {
        u.Height = *req.Height
    }
    if req.Weight != nil {
        u.Weight = *req.Weight
    }
    if req.Goal != nil {
        u.Goal = *req.Goal
    }
    if req.ActivityLevel != nil {
        u.ActivityLevel = *req.ActivityLevel
    }
    if req.Language != nil {
        u.Language = *req.Language
    }
    if req.NotificationsEnabled != nil {
        u.NotificationsEnabled = *req.NotificationsEnabled
    }
    // Пересчитываем BMI если изменился рост или вес
    u.CalculateBMI()
}