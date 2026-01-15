package handlers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/leroysb/go_kubernetes/internal/database"
	"github.com/leroysb/go_kubernetes/internal/database/models"
)

// GetOrders returns all orders with pagination
func GetOrders(c *fiber.Ctx) error {
	page := c.Query("page")
	if page == "" {
		page = "1"
	}

	// Convert page to integer
	pageNum, err := strconv.Atoi(page)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid page number"})
	}

	// Calculate offset and limit for pagination
	offset := (pageNum - 1) * 20
	limit := 20

	// Fetch orders from database in a goroutine
	var orders []models.Order
	done := make(chan bool)
	go func() {
		if err := database.DB.Db.Offset(offset).Limit(limit).Find(&orders).Where("status = ?", "paid").Error; err != nil {
			done <- false
		} else {
			done <- true
		}
	}()

	select {
	case success := <-done:
		if success {
			return c.Status(200).JSON(orders)
		} else {
			return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
		}
	case <-time.After(5 * time.Second):
		return c.Status(500).JSON(fiber.Map{"error": "Timeout"})
	}
}
