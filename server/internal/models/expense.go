package models

import (
	"time"
)

// SplitExpense represents how an expense is divided among group members
type SplitExpense struct {
	Base
	GroupID            string    `gorm:"type:uuid;not null;index" json:"groupId"`
	ExpenseID          string    `gorm:"type:uuid;not null;index" json:"expenseId"`
	CreatedBy          string    `gorm:"type:uuid;not null;index" json:"createdBy"`
	TotalAmount        float64   `gorm:"not null" json:"totalAmount"`
	SplitType          string    `gorm:"type:varchar(20);not null;default:'EQUAL';index" json:"splitType"`
	SettlementPriority int       `gorm:"default:0" json:"settlementPriority"`
	GraceEndDate       time.Time `gorm:"index;default:CURRENT_TIMESTAMP" json:"graceEndDate"`
	CustomSplitRules   []byte    `gorm:"type:jsonb;default:'{}'" json:"customSplitRules"`
	NeedsApproval      bool      `gorm:"default:false" json:"needsApproval"`
	DueDate            time.Time `gorm:"index;not null" json:"dueDate"`

	// Relations
	Group   Group        `gorm:"foreignKey:GroupID" json:"-"`
	Expense Expense      `gorm:"foreignKey:ExpenseID" json:"expense,omitempty"`
	Creator User         `gorm:"foreignKey:CreatedBy" json:"-"`
	Shares  []SplitShare `gorm:"foreignKey:SplitExpenseID" json:"shares,omitempty"`
}

// SplitShare represents an individual's portion of a split expense
type SplitShare struct {
	Base
	SplitExpenseID    string     `gorm:"type:uuid;not null;index" json:"splitExpenseId"`
	UserID            string     `gorm:"type:uuid;not null;index" json:"userId"`
	Amount            float64    `gorm:"type:decimal(10,2);not null" json:"amount"`
	IsPaid            bool       `gorm:"default:false" json:"isPaid"`
	PaidAt            *time.Time `json:"paidAt,omitempty"`
	InterestRate      float64    `gorm:"type:decimal(5,2);default:0" json:"interestRate"`
	InterestAccrued   float64    `gorm:"type:decimal(10,2);default:0" json:"interestAccrued"`
	NextReminderDate  *time.Time `json:"nextReminderDate,omitempty"`
	ReminderFrequency string     `gorm:"type:varchar(20);default:''" json:"reminderFrequency"`

	// Relations
	SplitExpense SplitExpense `gorm:"foreignKey:SplitExpenseID" json:"-"`
	User         User         `gorm:"foreignKey:UserID" json:"-"`
}

type Expense struct {
	Base
	UserID           string    `gorm:"type:uuid;not null;index" json:"userId"`
	GroupID          *string   `gorm:"type:uuid;index" json:"groupId"`
	Amount           float64   `gorm:"not null" json:"amount"`
	OriginalCurrency string    `gorm:"not null" json:"originalCurrency"`
	ConvertedAmount  float64   `json:"convertedAmount"`
	Category         string    `gorm:"not null" json:"category"`
	Description      string    `json:"description"`
	Date             time.Time `gorm:"not null" json:"date"`
	IsVerified       bool      `gorm:"default:false" json:"isVerified"`
	// Relations
	User          User           `gorm:"foreignKey:UserID" json:"-"`
	Group         *Group         `gorm:"foreignKey:GroupID" json:"group,omitempty"`
	SplitExpenses []SplitExpense `gorm:"foreignKey:ExpenseID" json:"splitExpenses,omitempty"`
}

// RecurringExpense represents a recurring expense pattern
type RecurringExpense struct {
	Base
	UserID        string     `gorm:"type:uuid;not null;index" json:"userId"`
	GroupID       *string    `gorm:"type:uuid;index" json:"groupId"`
	Amount        float64    `gorm:"not null" json:"amount"`
	Currency      string     `gorm:"not null" json:"currency"`
	Category      string     `gorm:"not null" json:"category"`
	Description   string     `json:"description"`
	Frequency     string     `gorm:"not null" json:"frequency"` // daily, weekly, monthly, yearly
	StartDate     time.Time  `gorm:"not null" json:"startDate"`
	EndDate       *time.Time `json:"endDate"`
	LastProcessed time.Time  `json:"lastProcessed"`
	NextDueDate   time.Time  `json:"nextDueDate"`
	IsAutomatic   bool       `gorm:"default:false" json:"isAutomatic"`
	ReminderDays  int        `gorm:"default:0" json:"reminderDays"`
	IsActive      bool       `gorm:"default:true" json:"isActive"`

	// Relations
	User  User   `gorm:"foreignKey:UserID" json:"-"`
	Group *Group `gorm:"foreignKey:GroupID" json:"group,omitempty"`
}

// SplitExpense represents how an expense is split among group members
