package middleware

import (
	"github.com/fadelm2/belajar_midtrans/config"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"time"
)

func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		requestID := uuid.New().String()
		c.Locals("request_id", requestID)

		err := c.Next()

		config.Log.WithFields(map[string]interface{}{
			"request_id": requestID,
			"method":     c.Method(),
			"path":       c.OriginalURL(),
			"status":     c.Response().StatusCode(),
			"latency":    time.Since(start).String(),
			"ip":         c.IP(),
		}).Info("request finished")

		return err
	}
}
