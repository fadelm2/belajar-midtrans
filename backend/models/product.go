package models

type Product struct {
	ID    uint `gorm:"primaryKey"`
	Name  string
	Price int
}
