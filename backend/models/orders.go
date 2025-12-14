package models

type Order struct {
	ID      uint `gorm:"primaryKey"`
	OrderID string
	Total   int
	Status  string
}
