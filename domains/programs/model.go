package programs

import "github.com/msskobelina/fit-profi/pkg/mysql"

// swagger:model trainingProgram
type TrainingProgram struct {
	ID          int          `json:"id,omitempty" gorm:"primaryKey"`
	UserID      int          `json:"userId"        gorm:"index;not null"`
	CoachUserID *int         `json:"coachUserId,omitempty" gorm:"index"`
	Title       string       `json:"title"`
	Notes       string       `json:"notes" gorm:"type:text"`
	Weeks       int          `json:"weeks"`
	VideoURL    *string      `json:"videoUrl,omitempty"`
	Viewed      bool         `json:"viewed" gorm:"default:false"`
	Days        []ProgramDay `gorm:"foreignKey:ProgramID;constraint:OnDelete:CASCADE"`

	mysql.Model
}

// swagger:model programDay
type ProgramDay struct {
	ID         int               `json:"id,omitempty" gorm:"primaryKey"`
	ProgramID  int               `json:"programId" gorm:"index"`
	WeekNumber int               `json:"weekNumber"`
	DayNumber  int               `json:"dayNumber"`
	Title      string            `json:"title"`
	Exercises  []ProgramExercise `gorm:"foreignKey:ProgramDayID;constraint:OnDelete:CASCADE"`

	mysql.Model
}

// swagger:model programExercise
type ProgramExercise struct {
	ID           int      `json:"id,omitempty" gorm:"primaryKey"`
	ProgramDayID int      `json:"programDayId" gorm:"index"`
	Name         string   `json:"name"`
	TargetSets   int      `json:"targetSets"`
	TargetReps   string   `json:"targetReps"`
	TargetWeight *float64 `json:"targetWeight,omitempty"`

	mysql.Model
}

// swagger:model exerciseProgress
type ExerciseProgress struct {
	ID                int      `json:"id,omitempty" gorm:"primaryKey"`
	ProgramExerciseID int      `json:"programExerciseId" gorm:"index"`
	WeekNumber        int      `json:"weekNumber"`
	SetNumber         int      `json:"setNumber"`
	Weight            *float64 `json:"weight,omitempty"`
	Reps              string   `json:"reps"`

	mysql.Model
}

// swagger:model createProgramRequest
type CreateProgramRequest struct {
	UserID int                `json:"userId"`
	Title  string             `json:"title"`
	Notes  string             `json:"notes"`
	Weeks  int                `json:"weeks"`
	Days   []CreateProgramDay `json:"days"`
}

type CreateProgramDay struct {
	WeekNumber int                     `json:"weekNumber"`
	DayNumber  int                     `json:"dayNumber"`
	Title      string                  `json:"title"`
	Exercises  []CreateProgramExercise `json:"exercises"`
}

type CreateProgramExercise struct {
	Name         string   `json:"name"`
	TargetSets   int      `json:"targetSets"`
	TargetReps   string   `json:"targetReps"`
	TargetWeight *float64 `json:"targetWeight,omitempty"`
}

// swagger:model addProgressRequest
type AddProgressRequest struct {
	ProgramExerciseID int      `json:"programExerciseId"`
	WeekNumber        int      `json:"weekNumber"`
	SetNumber         int      `json:"setNumber"`
	Weight            *float64 `json:"weight,omitempty"`
	Reps              string   `json:"reps"`
}

// swagger:parameters programsGet programsDelete
type ProgramIDPath struct {
	// in: path
	// required: true
	ID int `json:"id"`
}

// swagger:model errorResponse
type ErrorResponse struct {
	// example: message
	Error string `json:"error"`
}
