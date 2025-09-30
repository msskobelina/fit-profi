package integrations

import "github.com/msskobelina/fit-profi/pkg/mysql"

type Provider string

const ProviderGoogle Provider = "google"

// swagger:model userIntegration
type UserIntegration struct {
	ID           int      `json:"id,omitempty" gorm:"primaryKey"`
	UserID       int      `json:"userId" gorm:"index;not null"`
	Provider     Provider `json:"provider" gorm:"type:enum('google');not null"`
	AccessToken  string   `json:"-" gorm:"type:text;not null"`
	RefreshToken string   `json:"-" gorm:"type:text"`
	ExpiryUnix   int64    `json:"expiryUnix" gorm:"index"`
	Scope        string   `json:"scope" gorm:"type:text"`
	CalendarID   string   `json:"calendarId" gorm:"type:varchar(256);default:'primary'"`
	Timezone     string   `json:"timezone" gorm:"type:varchar(64)"`
	mysql.Model
}
