package models

import "time"

// SplitExpense represents how an expense is divided among group members
type SplitExpense struct {
	Base
	GroupID            string    `gorm:"type:uuid;not null;index" json:"groupId"`
	ExpenseID          string    `gorm:"type:uuid;not null;index" json:"expenseId"`
	CreatedBy          string    `gorm:"type:uuid;not null;index" json:"createdBy"`
	TotalAmount        float64   `gorm:"not null" json:"totalAmount"`
	SplitType          SplitType `gorm:"type:varchar(20);not null;index" json:"splitType"`
	SettlementPriority int       `json:"settlementPriority"`
	GraceEndDate       time.Time `gorm:"index" json:"graceEndDate"`
	CustomSplitRules   JSON      `gorm:"type:jsonb" json:"customSplitRules"`
	NeedsApproval      bool      `json:"needsApproval"`
	DueDate            time.Time `gorm:"index" json:"dueDate"`

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
	Amount            float64    `gorm:"not null" json:"amount"`
	IsPaid            bool       `gorm:"index" json:"isPaid"`
	PaidAt            *time.Time `json:"paidAt"`
	InterestRate      float64    `json:"interestRate"`
	InterestAccrued   float64    `json:"interestAccrued"`
	NextReminderDate  *time.Time `gorm:"index" json:"nextReminderDate"`
	ReminderFrequency string     `json:"reminderFrequency"`

	// Relations
	SplitExpense SplitExpense `gorm:"foreignKey:SplitExpenseID" json:"-"`
	User         User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Payments     []Payment    `gorm:"foreignKey:SplitShareID" json:"payments,omitempty"`
}

type Expense struct {
	Base
	UserID           string    `gorm:"type:uuid;not null;index" json:"userId"`
	GroupID          *string   `gorm:"type:uuid;index" json:"groupId"`
	Amount           float64   `gorm:"not null" json:"amount"`
	OriginalCurrency string    `gorm:"not null" json:"originalCurrency"`
	ConvertedAmount  float64   `json:"convertedAmount"`
	Category         string    `gorm:"not null" json:"category"`
	Tags             []string  `gorm:"type:text[]" json:"tags"`
	Description      string    `json:"description"`
	Date             time.Time `gorm:"not null" json:"date"`
	IsVerified       bool      `gorm:"default:false" json:"isVerified"`
	ImageURL         string    `json:"imageUrl"`
	EmotionalState   string    `json:"emotionalState"`
	IsImpulsive      bool      `gorm:"default:false" json:"isImpulsive"`
	MindfulnessScore float64   `gorm:"default:0" json:"mindfulnessScore"`
	IsRecurring      bool      `gorm:"default:false" json:"isRecurring"`
	IsEssential      bool      `gorm:"default:false" json:"isEssential"`
	PaymentMode      string    `json:"paymentMode"`

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
