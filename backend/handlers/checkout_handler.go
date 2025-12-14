package handlers

import (
	"fmt"
	"time"

	"github.com/fadelm2/belajar_midtrans/config"
	"github.com/fadelm2/belajar_midtrans/models"

	"github.com/gofiber/fiber/v2"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type CheckoutRequest struct {
	Items []struct {
		ProductID uint `json:"product_id"`
		Qty       int  `json:"qty"`
	} `json:"items"`
}

func Checkout(c *fiber.Ctx) error {
	var req CheckoutRequest
	if err := c.BodyParser(&req); err != nil || len(req.Items) == 0 {
		return fiber.NewError(400, "invalid request")
	}

	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	orderID := fmt.Sprintf("ORDER-%d", time.Now().Unix())
	expiresAt := time.Now().Add(5 * time.Minute)
	now := time.Now()

	// =========================
	// CREATE ORDER
	// =========================
	order := models.Order{
		OrderID:         orderID,
		Status:          "PENDING",
		ExpiresAt:       expiresAt,
		TransactionTime: now,
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return fiber.NewError(500, "failed create order")
	}

	// =========================
	// COLLECT PRODUCT IDS
	// =========================
	productIDs := []uint{}
	qtyMap := map[uint]int{}

	for _, i := range req.Items {
		productIDs = append(productIDs, i.ProductID)
		qtyMap[i.ProductID] = i.Qty
	}

	// =========================
	// FETCH PRODUCTS
	// =========================
	var products []models.Product
	if err := tx.Where("id IN ?", productIDs).Find(&products).Error; err != nil {
		tx.Rollback()
		return fiber.NewError(500, "failed fetch products")
	}

	if len(products) != len(productIDs) {
		tx.Rollback()
		return fiber.NewError(404, "product not found")
	}

	// =========================
	// CREATE ORDER ITEMS
	// =========================
	var items []models.OrderItem
	for _, p := range products {
		items = append(items, models.OrderItem{
			OrderID:   orderID,
			ProductID: p.ID,
			Quantity:  qtyMap[p.ID],
			Price:     p.Price,
		})
	}

	if err := tx.Create(&items).Error; err != nil {
		tx.Rollback()
		return fiber.NewError(500, "failed create order items")
	}

	// =========================
	// AGGREGATE TOTAL
	// =========================
	var total int
	if err := tx.Raw(`
		SELECT COALESCE(SUM(price * quantity), 0)
		FROM order_items
		WHERE order_id = ?
	`, orderID).Scan(&total).Error; err != nil {
		tx.Rollback()
		return fiber.NewError(500, "failed calculate total")
	}

	if err := tx.Model(&models.Order{}).
		Where("order_id = ?", orderID).
		Update("total", total).Error; err != nil {
		tx.Rollback()
		return fiber.NewError(500, "failed update total")
	}

	if err := tx.Commit().Error; err != nil {
		return fiber.NewError(500, "commit failed")
	}

	// =========================
	// PAYMENT RULE
	// =========================
	enabledPayments := []snap.SnapPaymentType{
		snap.PaymentTypeGopay,
	}

	if total >= 110000 {
		enabledPayments = append(enabledPayments, snap.PaymentTypeBankTransfer)
	}

	// =========================
	// MIDTRANS
	// =========================
	s := snap.Client{}
	s.New("SB-Mid-server-YyzkaYvpKpLec8Cexdrb3LH7", midtrans.Sandbox)

	resp, err := s.CreateTransaction(&snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: int64(total),
		},
		EnabledPayments: enabledPayments,
	})
	if err != nil {
		return fiber.NewError(500, "midtrans error")
	}

	return c.JSON(fiber.Map{
		"order_id":   orderID,
		"total":      total,
		"expires_at": expiresAt,
		"token":      resp.Token,
	})
}
