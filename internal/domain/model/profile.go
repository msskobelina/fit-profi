package model

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
