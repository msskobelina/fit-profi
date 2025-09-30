package nutrition

import (
	"time"

	"github.com/msskobelina/fit-profi/pkg/mysql"
)

// swagger:model mealType
type MealType string

const (
	MealBreakfast MealType = "breakfast"
	MealLunch     MealType = "lunch"
	MealDinner    MealType = "dinner"
	MealSnack     MealType = "snack"
)

// swagger:model diaryEntry
type DiaryEntry struct {
	ID     int       `json:"id,omitempty" gorm:"primaryKey"`
	UserID int       `json:"userId"        gorm:"index;not null"`
	Date   time.Time `json:"date"          gorm:"type:date;index;not null"`
	Meal   MealType  `json:"meal"          gorm:"type:enum('breakfast','lunch','dinner','snack');not null"`
	Notes  string    `json:"notes"         gorm:"type:text"`

	TotalCalories float32 `json:"totalCalories"`
	TotalProtein  float32 `json:"totalProtein"`
	TotalFat      float32 `json:"totalFat"`
	TotalCarbs    float32 `json:"totalCarbs"`

	Items []DiaryItem `json:"items" gorm:"constraint:OnDelete:CASCADE"`

	mysql.Model
}

// swagger:model diaryItem
type DiaryItem struct {
	ID           int     `json:"id,omitempty" gorm:"primaryKey"`
	DiaryEntryID int     `json:"diaryEntryId" gorm:"index;not null"`
	Name         string  `json:"name"         gorm:"not null"`
	Quantity     float32 `json:"quantity"`
	Unit         string  `json:"unit"`
	Calories     float32 `json:"calories"`
	Protein      float32 `json:"protein"`
	Fat          float32 `json:"fat"`
	Carbs        float32 `json:"carbs"`

	mysql.Model
}

// swagger:model createDiaryItem
type CreateDiaryItem struct {
	Name     string  `json:"name"`
	Quantity float32 `json:"quantity"`
	Unit     string  `json:"unit"`
	Calories float32 `json:"calories"`
	Protein  float32 `json:"protein"`
	Fat      float32 `json:"fat"`
	Carbs    float32 `json:"carbs"`
}

// swagger:model createDiaryEntryRequest
type CreateDiaryEntryRequest struct {
	Date string `json:"date"`
	// enum: breakfast,lunch,dinner,snack
	Meal  MealType          `json:"meal"`
	Notes string            `json:"notes"`
	Items []CreateDiaryItem `json:"items"`
}

// swagger:model updateDiaryEntryRequest
type UpdateDiaryEntryRequest struct {
	Date  string            `json:"date"`
	Meal  MealType          `json:"meal"`
	Notes string            `json:"notes"`
	Items []CreateDiaryItem `json:"items"`
}

// swagger:model diarySummary
type DiarySummary struct {
	From          string  `json:"from"`
	To            string  `json:"to"`
	TotalCalories float32 `json:"totalCalories"`
	TotalProtein  float32 `json:"totalProtein"`
	TotalFat      float32 `json:"totalFat"`
	TotalCarbs    float32 `json:"totalCarbs"`
}

// swagger:model errorResponse
type ErrorResponse struct {
	Error string `json:"error"`
}
