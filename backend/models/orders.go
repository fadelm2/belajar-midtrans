package models

import "time"

type Order struct {
	ID              uint      `gorm:"primaryKey;autoIncrement"`
	OrderID         string    `gorm:"column:order_id;type:varchar(50);uniqueIndex"`
	Total           int       `gorm:"type:int"`
	Status          string    `gorm:"type:varchar(20)"`
	PaymentType     string    `gorm:"column:payment_type;type:varchar(50)"`
	TransactionTime time.Time `gorm:"column:transaction_time;type:datetime"`
	ExpiresAt       time.Time `gorm:"column:expires_at;type:datetime"`
}
