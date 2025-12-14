package handlers

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"

	"github.com/fadelm2/belajar_midtrans/config"
	"github.com/fadelm2/belajar_midtrans/models"

	"github.com/gofiber/fiber/v2"
)

type MidtransWebhook struct {
	OrderID           string `json:"order_id"`
	TransactionStatus string `json:"transaction_status"`
	GrossAmount       string `json:"gross_amount"`
	SignatureKey      string `json:"signature_key"`
	PaymentType       string `json:"payment_type"`
	TransactionTime   string `json:"transaction_time"`
	StatusCode        string `json:"status_code"`
}

func MidtransWebhookHandler(c *fiber.Ctx) error {
	var payload MidtransWebhook
	if err := c.BodyParser(&payload); err != nil {
		return fiber.NewError(400, "invalid payload")
	}

	// =========================
	// 1️⃣ VALIDASI SIGNATURE
	// =========================
	raw := payload.OrderID +
		payload.StatusCode +
		payload.GrossAmount +
		"SB-Mid-server-YyzkaYvpKpLec8Cexdrb3LH7"

	hash := sha512.Sum512([]byte(raw))
	expectedSignature := hex.EncodeToString(hash[:])

	if expectedSignature != payload.SignatureKey {
		return fiber.NewError(401, "invalid signature")
	}

	// =========================
	// 2️⃣ AMBIL ORDER
	// =========================
	var order models.Order
	if err := config.DB.
		Where("order_id = ?", payload.OrderID).
		First(&order).Error; err != nil {
		return fiber.NewError(404, "order not found")
	}

	// =========================
	// 3️⃣ VALIDASI AMOUNT
	// =========================
	expectedAmount := fmt.Sprintf("%d.00", order.Total)
	if payload.GrossAmount != expectedAmount {
		return fiber.NewError(400, "amount mismatch")
	}

	// =========================
	// 4️⃣ IDEMPOTENT CHECK
	// =========================
	if order.Status == "PAID" {
		return c.JSON(fiber.Map{
			"message": "already processed",
		})
	}

	// =========================
	// 5️⃣ MAP STATUS
	// =========================
	newStatus := "PENDING"

	switch payload.TransactionStatus {
	case "settlement", "capture":
		newStatus = "PAID"
	case "pending":
		newStatus = "PENDING"
	case "expire":
		newStatus = "EXPIRED"
	case "cancel":
		newStatus = "CANCELLED"
	case "deny":
		newStatus = "FAILED"
	}

	// =========================
	// 6️⃣ UPDATE ORDER
	// =========================
	config.DB.Model(&models.Order{}).
		Where("order_id = ?", payload.OrderID).
		Updates(map[string]interface{}{
			"status":           newStatus,
			"payment_type":     payload.PaymentType,
			"transaction_time": payload.TransactionTime,
		})

	return c.JSON(fiber.Map{
		"message": "ok",
	})
}
