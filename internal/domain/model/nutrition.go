package model

import (
	"time"

	"github.com/msskobelina/fit-profi/pkg/mysql"
)

type DiaryEntry struct {
	ID       int         `json:"id,omitempty" gorm:"primaryKey"`
	UserID   int         `json:"userId" gorm:"index;not null"`
	Date     time.Time   `json:"date" gorm:"index;not null"`
	MealType string      `json:"mealType" gorm:"type:enum('breakfast','lunch','dinner','snack');not null"`
	Items    []DiaryItem `json:"items,omitempty" gorm:"foreignKey:EntryID;constraint:OnDelete:CASCADE"`

	mysql.Model
}

type DiaryItem struct {
	ID       int     `json:"id,omitempty" gorm:"primaryKey"`
	EntryID  int     `json:"entryId" gorm:"index;not null"`
	Name     string  `json:"name"`
	Grams    float32 `json:"grams"`
	Calories float32 `json:"calories"`
	ProteinG float32 `json:"proteinG"`
	FatG     float32 `json:"fatG"`
	CarbsG   float32 `json:"carbsG"`

	mysql.Model
}
