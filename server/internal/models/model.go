package models

import (
	"time"

	"gorm.io/gorm"
)

type Base struct {
	ID        string         `gorm:"primarykey;type:uuid;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type JSON map[string]interface{}

// Enum types
type GroupType string
type SplitType string

const (
	GroupTypeCouple   GroupType = "COUPLE"
	GroupTypeFlatmate GroupType = "FLATMATE"
	GroupTypeTrip     GroupType = "TRIP"
	GroupTypeHome     GroupType = "HOME"
	GroupTypeOther    GroupType = "OTHER"

	SplitTypeEqual      SplitType = "EQUAL"
	SplitTypePercentage SplitType = "PERCENTAGE"
	SplitTypeCustom     SplitType = "CUSTOM"
	SplitTypeShares     SplitType = "SHARES"
)
