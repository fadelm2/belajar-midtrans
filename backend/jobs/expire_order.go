package jobs

import (
	"github.com/fadelm2/belajar_midtrans/config"
	"github.com/fadelm2/belajar_midtrans/models"
	"time"
)

func StartExpireOrderJob() {
	ticker := time.NewTicker(1 * time.Minute)

	go func() {
		for range ticker.C {
			config.DB.Model(&models.Order{}).
				Where("status = ? AND expires_at < ?", "PENDING", time.Now()).
				Update("status", "EXPIRED")
		}
	}()
}
