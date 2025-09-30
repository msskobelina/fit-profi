package calendar

import (
	"time"

	"github.com/msskobelina/fit-profi/pkg/mysql"
)

// swagger:model coachAvailability
type CoachAvailability struct {
	ID        int       `json:"id,omitempty" gorm:"primaryKey"`
	CoachID   int       `json:"coachId" gorm:"index;not null"`
	StartTime time.Time `json:"startTime" gorm:"index;not null"`
	EndTime   time.Time `json:"endTime"   gorm:"not null"`
	Notes     string    `json:"notes" gorm:"type:text"`
	mysql.Model
}
