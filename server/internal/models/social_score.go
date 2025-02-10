package models

import (
	"time"
)

type SocialScoreHistory struct {
	Base
	UserID    string    `gorm:"type:uuid;not null;index" json:"userId"`
	OldScore  float64   `json:"oldScore"`
	NewScore  float64   `json:"newScore"`
	Reason    string    `json:"reason"`
	Timestamp time.Time `gorm:"index" json:"timestamp"`

	// Relations
	User User `gorm:"foreignKey:UserID" json:"-"`
}
