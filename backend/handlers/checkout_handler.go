package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/midtrans/midtrans-go"
	"time"

	"github.com/fadelm2/belajar_midtrans/config"
	"github.com/fadelm2/belajar_midtrans/models"

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
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(400, "invalid request")
	}

	tx := config.DB.Begin()

	orderID := fmt.Sprintf("ORDER-%d", time.Now().Unix())

	// create order
	tx.Create(&models.Order{
		OrderID: orderID,
		Status:  "PENDING",
	})

	// collect product ids
	productIDs := []uint{}
	qtyMap := map[uint]int{}

	for _, i := range req.Items {
		productIDs = append(productIDs, i.ProductID)
		qtyMap[i.ProductID] = i.Qty
	}

	// fetch products ONCE
	var products []models.Product
	tx.Where("id IN ?", productIDs).Find(&products)

	if len(products) != len(productIDs) {
		tx.Rollback()
		return fiber.NewError(404, "product not found")
	}

	// build order items
	items := []models.OrderItem{}
	for _, p := range products {
		items = append(items, models.OrderItem{
			OrderID:   orderID,
			ProductID: p.ID,
			Quantity:  qtyMap[p.ID],
			Price:     p.Price,
		})
	}

	tx.Create(&items)

	// aggregate total
	type Result struct {
		Total int
	}
	var result Result
	tx.Raw(`
		SELECT SUM(price * quantity) AS total
		FROM order_items
		WHERE order_id = ?
	`, orderID).Scan(&result)

	tx.Model(&models.Order{}).
		Where("order_id = ?", orderID).
		Update("total", result.Total)

	tx.Commit()

	// midtrans
	s := snap.Client{}
	s.New("SB-Mid-server-YyzkaYvpKpLec8Cexdrb3LH7", midtrans.Sandbox)

	resp, err := s.CreateTransaction(&snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: int64(result.Total),
		},
	})
	if err != nil {
		return fiber.NewError(500, "midtrans error")
	}

	return c.JSON(fiber.Map{
		"order_id": orderID,
		"total":    result.Total,
		"token":    resp.Token,
	})
}
