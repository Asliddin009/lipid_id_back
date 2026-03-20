package models

// FoodRecognitionResult — результат распознавания еды по фото
type FoodRecognitionResult struct {
	DishName      string  `json:"dishName"`
	Calories      int     `json:"calories"`
	Proteins      float64 `json:"proteins"`
	Fats          float64 `json:"fats"`
	SaturatedFats float64 `json:"saturatedFats"`
	Carbs         float64 `json:"carbs"`
	LipidRating   string  `json:"lipidRating"`
	RatingImpact  float64 `json:"ratingImpact"`
	Confidence    float64 `json:"confidence"`
}
