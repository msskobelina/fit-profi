package model

import "github.com/msskobelina/fit-profi/pkg/mysql"

type TrainingProgram struct {
	ID          int          `json:"id,omitempty" gorm:"primaryKey"`
	UserID      int          `json:"userId" gorm:"index;not null"`
	Title       string       `json:"title"`
	Description string       `json:"description" gorm:"type:text"`
	Days        []ProgramDay `json:"days,omitempty" gorm:"foreignKey:ProgramID;constraint:OnDelete:CASCADE"`

	mysql.Model
}

type ProgramDay struct {
	ID        int               `json:"id,omitempty" gorm:"primaryKey"`
	ProgramID int               `json:"programId" gorm:"index;not null"`
	DayNumber int               `json:"dayNumber"`
	Title     string            `json:"title"`
	Exercises []ProgramExercise `json:"exercises,omitempty" gorm:"foreignKey:DayID;constraint:OnDelete:CASCADE"`

	mysql.Model
}

type ProgramExercise struct {
	ID       int    `json:"id,omitempty" gorm:"primaryKey"`
	DayID    int    `json:"dayId" gorm:"index;not null"`
	Name     string `json:"name"`
	Sets     int    `json:"sets"`
	Reps     int    `json:"reps"`
	WeightKg int    `json:"weightKg"`
	Notes    string `json:"notes" gorm:"type:text"`

	mysql.Model
}

type ExerciseProgress struct {
	ID         int    `json:"id,omitempty" gorm:"primaryKey"`
	UserID     int    `json:"userId" gorm:"index;not null"`
	ExerciseID int    `json:"exerciseId" gorm:"index;not null"`
	Sets       int    `json:"sets"`
	Reps       int    `json:"reps"`
	WeightKg   int    `json:"weightKg"`
	Notes      string `json:"notes" gorm:"type:text"`

	mysql.Model
}
