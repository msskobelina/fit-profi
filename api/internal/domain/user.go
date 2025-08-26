package domain

import "github.com/msskobelina/fit-profi/pkg/mysql"

type User struct {
	ID       int    `json:"id,omitempty" gorm:"primaryKey"`
	FullName string `json:"fullName,omitempty"`
	Email    string `json:"email,omitempty" gorm:"not null;unique;index"`
	Password string `json:"password,omitempty"`

	mysql.Model
}

type UserToken struct {
	ID    int    `json:"id,omitempty" gorm:"primaryKey"`
	Email string `json:"email,omitempty" gorm:"index"`
	Token string `json:"token,omitempty"`

	mysql.Model
}

type RevokedToken struct {
	ID        int    `gorm:"primaryKey"`
	JTI       string `gorm:"uniqueIndex;size:64"`
	ExpiresAt int64

	mysql.Model
}
