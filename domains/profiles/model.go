package profiles

import (
	"time"

	"github.com/msskobelina/fit-profi/pkg/mysql"
)

// swagger:model userProfile
type UserProfile struct {
	ID          int     `json:"id,omitempty" gorm:"primaryKey"`
	UserID      int     `json:"userId"        gorm:"uniqueIndex;not null"`
	FullName    string  `json:"fullName"`
	Age         int     `json:"age"`
	WeightKg    float32 `json:"weightKg"`
	Goal        Goal    `json:"goal" gorm:"type:enum('lose_weight','gain_weight','rehab','keep_fit','competition')"`
	Description string  `json:"description" gorm:"type:text"`

	mysql.Model
}

type Goal string

const (
	GoalLoseWeight  Goal = "lose_weight"
	GoalGainWeight  Goal = "gain_weight"
	GoalRehab       Goal = "rehab"
	GoalKeepFit     Goal = "keep_fit"
	GoalCompetition Goal = "competition"
)

// swagger:model coachProfile
type CoachProfile struct {
	ID           int                `json:"id,omitempty" gorm:"primaryKey"`
	UserID       int                `json:"userId"        gorm:"uniqueIndex;not null"`
	FullName     string             `json:"fullName"`
	Category     CoachCategory      `json:"category" gorm:"type:enum('standard','master','professional')"`
	Info         string             `json:"info"     gorm:"type:text"`
	Achievements []CoachAchievement `json:"achievements" gorm:"constraint:OnDelete:CASCADE"`
	Education    []CoachEducation   `json:"education"    gorm:"constraint:OnDelete:CASCADE"`

	mysql.Model
}

type CoachCategory string

const (
	CoachStandard CoachCategory = "standard"
	CoachMaster   CoachCategory = "master"
	CoachPro      CoachCategory = "professional"
)

// swagger:model coachAchievement
type CoachAchievement struct {
	ID             int       `json:"id,omitempty" gorm:"primaryKey"`
	CoachProfileID int       `json:"coachProfileId" gorm:"index;not null"`
	StartPeriod    time.Time `json:"startPeriod"`
	EndPeriod      time.Time `json:"endPeriod"`
	Title          string    `json:"title"`
	CertificateURL *string   `json:"certificateUrl,omitempty"`

	mysql.Model
}

// swagger:model coachEducation
type CoachEducation struct {
	ID             int       `json:"id,omitempty" gorm:"primaryKey"`
	CoachProfileID int       `json:"coachProfileId" gorm:"index;not null"`
	StartPeriod    time.Time `json:"startPeriod"`
	EndPeriod      time.Time `json:"endPeriod"`
	Place          string    `json:"place"`
	Description    string    `json:"description,omitempty"`

	mysql.Model
}

// swagger:model createUserProfileRequest
type CreateUserProfileRequest struct {
	FullName    string  `json:"fullName"`
	Age         int     `json:"age"`
	WeightKg    float32 `json:"weightKg"`
	Goal        Goal    `json:"goal"`
	Description string  `json:"description"`
}

// swagger:model updateUserProfileRequest
type UpdateUserProfileRequest struct {
	FullName    *string  `json:"fullName,omitempty"`
	Age         *int     `json:"age,omitempty"`
	WeightKg    *float32 `json:"weightKg,omitempty"`
	Goal        *Goal    `json:"goal,omitempty"`
	Description *string  `json:"description,omitempty"`
}

// swagger:model createCoachProfileRequest
type CreateCoachProfileRequest struct {
	FullName     string                   `json:"fullName"`
	Category     CoachCategory            `json:"category"`
	Info         string                   `json:"info"`
	Achievements []CreateCoachAchievement `json:"achievements"`
	Education    []CreateCoachEducation   `json:"education"`
}

type CreateCoachAchievement struct {
	CoachProfileID int       `json:"coachProfileId" gorm:"index;not null"`
	StartPeriod    time.Time `json:"startPeriod"`
	EndPeriod      time.Time `json:"endPeriod"`
	Title          string    `json:"title"`
	CertificateURL *string   `json:"certificateUrl,omitempty"`
}
type CreateCoachEducation struct {
	CoachProfileID int       `json:"coachProfileId" gorm:"index;not null"`
	StartPeriod    time.Time `json:"startPeriod"`
	EndPeriod      time.Time `json:"endPeriod"`
	Place          string    `json:"place"`
	Description    string    `json:"description"`
}

// swagger:model updateCoachProfileRequest
type UpdateCoachProfileRequest struct {
	FullName     *string                   `json:"fullName,omitempty"`
	Category     *CoachCategory            `json:"category,omitempty"`
	Info         *string                   `json:"info,omitempty"`
	Achievements *[]CreateCoachAchievement `json:"achievements,omitempty"`
	Education    *[]CreateCoachEducation   `json:"education,omitempty"`
}

// swagger:parameters profilesAdminGetUser
type GetUserProfilePath struct {
	// in: path
	// required: true
	UserID int `json:"userId"`
}

// swagger:model errorResponse
type ErrorResponse struct {
	Error string `json:"error"`
}
