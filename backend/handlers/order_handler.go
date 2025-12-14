package handlers

import (
	"errors"
	"github.com/fadelm2/belajar_midtrans/config"
	"github.com/fadelm2/belajar_midtrans/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetOrderStatus(c *fiber.Ctx) error {
	orderID := c.Params("orderId")

	var order models.Order
	err := config.DB.
		Where("order_id = ?", orderID).
		First(&order).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{
				"message":  "Order belum dibuat atau belum tersimpan",
				"order_id": orderID,
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	return c.JSON(order)
}
