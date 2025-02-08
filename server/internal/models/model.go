package models

import (
	"database/sql/driver"
	"errors"
	"time"

	"gorm.io/gorm"
)

type Base struct {
	ID        string         `gorm:"primarykey;type:uuid;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type JSON []byte

// Scan implements the sql.Scanner interface
func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	s, ok := value.([]byte)
	if !ok {
		return errors.New("invalid scan source for JSON")
	}
	*j = append((*j)[0:0], s...)
	return nil
}

// Value implements the driver.Valuer interface
func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return string(j), nil
}

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
