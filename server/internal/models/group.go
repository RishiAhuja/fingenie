package models

import "time"

// Group represents a group of users who share expenses
type Group struct {
	Base
	Name                    string  `gorm:"not null" json:"name"`
	CreatedBy               string  `gorm:"type:uuid;not null" json:"createdBy"`
	Description             string  `json:"description"`
	DefaultCurrency         string  `gorm:"not null;default:'USD'" json:"defaultCurrency"`
	GroupType               string  `gorm:"not null" json:"groupType"` // household, rental, relationship, custom
	IsRecurring             bool    `gorm:"default:false" json:"isRecurring"`
	SecurityDepositRequired float64 `gorm:"default:0" json:"securityDepositRequired"`
	RequiresAdminApproval   bool    `gorm:"default:false" json:"requiresAdminApproval"`

	BudgetStrategy string `json:"budgetStrategy"`

	BillingCycleStart time.Time `json:"billingCycleStart"`
	SplitStrategy     string    `gorm:"not null;default:'equal'" json:"splitStrategy"`
	AutoSettlement    bool      `gorm:"default:false" json:"autoSettlement"`

	// Relations
	Creator           User               `gorm:"foreignKey:CreatedBy" json:"-"`
	Members           []GroupMember      `gorm:"foreignKey:GroupID" json:"members,omitempty"`
	Expenses          []Expense          `gorm:"foreignKey:GroupID" json:"expenses,omitempty"`
	RecurringExpenses []RecurringExpense `gorm:"foreignKey:GroupID" json:"recurringExpenses,omitempty"`
	Budgets           []Budget           `gorm:"foreignKey:GroupID" json:"budgets,omitempty"`
}

// GroupMember represents a user's membership in a group
type GroupMember struct {
	Base
	GroupID      string    `gorm:"type:uuid;not null;index" json:"groupId"`
	UserID       string    `gorm:"type:uuid;not null;index" json:"userId"`
	Role         string    `gorm:"not null;default:'member'" json:"role"` // admin, member
	JoinedAt     time.Time `gorm:"not null" json:"joinedAt"`
	IsActive     bool      `gorm:"default:true" json:"isActive"`
	SharePercent float64   `gorm:"default:0" json:"sharePercent"`

	// Relations
	Group Group `gorm:"foreignKey:GroupID" json:"-"`
	User  User  `gorm:"foreignKey:UserID" json:"-"`
}
