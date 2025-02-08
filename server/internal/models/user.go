package models

import (
	"time"
)

type User struct {
	Base
	DisplayName            string     `gorm:"not null" json:"displayName"`
	Email                  string     `gorm:"uniqueIndex;not null" json:"email"`
	Password               string     `gorm:"not null" json:"-"`
	PhoneNumber            string     `gorm:"uniqueIndex" json:"phoneNumber"`
	SocialScore            float64    `gorm:"default:0" json:"socialScore"`
	IsPremium              bool       `gorm:"default:false" json:"isPremium"`
	PreferredCurrency      string     `gorm:"not null;default:'USD'" json:"preferredCurrency"`
	TelegramID             string     `json:"telegramId"`
	WhatsappNumber         string     `json:"whatsappNumber"`
	CreditLimit            float64    `gorm:"default:0" json:"creditLimit"`
	NextSalaryDate         *time.Time `json:"nextSalaryDate"`
	HasDefaultHistory      bool       `gorm:"default:false" json:"hasDefaultHistory"`
	SecurityDepositBalance float64    `gorm:"default:0" json:"securityDepositBalance"`

	// Relations
	Expenses           []Expense            `gorm:"foreignKey:UserID" json:"expenses,omitempty"`
	GroupMembers       []GroupMember        `gorm:"foreignKey:UserID" json:"groupMembers,omitempty"`
	Budgets            []Budget             `gorm:"foreignKey:UserID" json:"budgets,omitempty"`
	IncomeStreams      []IncomeStream       `gorm:"foreignKey:UserID" json:"incomeStreams,omitempty"`
	BehavioralPatterns []BehavioralPattern  `gorm:"foreignKey:UserID" json:"behavioralPatterns,omitempty"`
	ScoreHistory       []SocialScoreHistory `gorm:"foreignKey:UserID" json:"scoreHistory,omitempty"`
}

type IncomeStream struct {
	Base
	UserID       string     `gorm:"type:uuid;not null;index" json:"userId"`
	Name         string     `gorm:"not null" json:"name"`
	Type         string     `gorm:"not null" json:"type"`
	Amount       float64    `gorm:"not null" json:"amount"`
	Frequency    string     `gorm:"not null" json:"frequency"`
	LastReceived *time.Time `json:"lastReceived"`
	NextExpected *time.Time `json:"nextExpected"`
	TaxCategory  string     `json:"taxCategory"`
	IsFreelance  bool       `gorm:"default:false" json:"isFreelance"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}

type BehavioralPattern struct {
	Base
	UserID         string   `gorm:"type:uuid;not null;index" json:"userId"`
	PatternType    string   `gorm:"not null" json:"patternType"`
	Frequency      float64  `json:"frequency"`
	TriggerFactors []string `gorm:"type:text[]" json:"triggerFactors"`
	RiskLevel      string   `json:"riskLevel"`
	Suggestions    JSON     `gorm:"type:jsonb" json:"suggestions"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}

type Budget struct {
	Base
	UserID            string    `gorm:"type:uuid;not null;index" json:"userId"`
	GroupID           *string   `gorm:"type:uuid;index" json:"groupId"`
	Category          string    `gorm:"not null" json:"category"`
	Tags              []string  `gorm:"type:text[];serializer:json" json:"tags"`
	Amount            float64   `gorm:"not null" json:"amount"`
	Period            string    `gorm:"not null" json:"period"`
	StartDate         time.Time `gorm:"not null" json:"startDate"`
	EndDate           time.Time `gorm:"not null" json:"endDate"`
	CurrentSpent      float64   `gorm:"default:0" json:"currentSpent"`
	AISuggestedAmount float64   `json:"aiSuggestedAmount"`
	IsAutoAdjusting   bool      `gorm:"default:false" json:"isAutoAdjusting"`

	User  User   `gorm:"foreignKey:UserID" json:"-"`
	Group *Group `gorm:"foreignKey:GroupID" json:"group,omitempty"`
}
