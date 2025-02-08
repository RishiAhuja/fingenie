package models

type Payment struct {
	Base
	SplitShareID      string  `gorm:"type:uuid;not null;index" json:"splitShareId"`
	FromUserID        string  `gorm:"type:uuid;not null;index" json:"fromUserId"`
	ToUserID          string  `gorm:"type:uuid;not null;index" json:"toUserId"`
	Amount            float64 `gorm:"not null" json:"amount"`
	Currency          string  `gorm:"not null" json:"currency"`
	PaymentMethod     string  `json:"paymentMethod"`
	Status            string  `gorm:"index" json:"status"`
	TransactionID     string  `json:"transactionId"`
	IsSecurityDeposit bool    `json:"isSecurityDeposit"`

	SplitShare SplitShare `gorm:"foreignKey:SplitShareID" json:"splitShare,omitempty"`
	FromUser   User       `gorm:"foreignKey:FromUserID" json:"fromUser,omitempty"`
	ToUser     User       `gorm:"foreignKey:ToUserID" json:"toUser,omitempty"`
}
