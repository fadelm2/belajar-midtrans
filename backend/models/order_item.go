package models

type OrderItem struct {
	ID        uint `gorm:"primaryKey"`
	OrderID   string
	ProductID uint
	Quantity  int
	Price     int
}
