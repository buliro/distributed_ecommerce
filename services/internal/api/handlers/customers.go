package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/leroysb/go_kubernetes/internal/api/auth"
	"github.com/leroysb/go_kubernetes/internal/database"
	"github.com/leroysb/go_kubernetes/internal/database/models"
	"github.com/leroysb/go_kubernetes/internal/sms"
	"github.com/leroysb/go_kubernetes/internal/utils"
	"gorm.io/gorm"
)

// CreateCustomer creates a new customer
func CreateCustomer(c *fiber.Ctx) error {
	customer := new(models.Customer)

	// Error check fields
	if err := c.BodyParser(customer); err != nil {
		if _, ok := err.(*json.UnmarshalTypeError); ok {
			if strings.Contains(err.Error(), "name") {
				return c.Status(400).JSON(fiber.Map{"error": "Missing name of type string"})
			}
			if strings.Contains(err.Error(), "phone") {
				return c.Status(400).JSON(fiber.Map{"error": "Missing phone of type string"})
			}
			if strings.Contains(err.Error(), "password") {
				return c.Status(400).JSON(fiber.Map{"error": "Missing password of type string"})
			}
		}
		if strings.Contains(err.Error(), "unexpected end of JSON input") {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON input"})
		}
		if errors.Is(err, fiber.ErrUnprocessableEntity) {
			return c.Status(422).JSON(fiber.Map{"error": "Unprocessable Entity"})
		}
		return c.Status(400).SendString(err.Error())
	}

	if customer.Name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing name"})
	}

	if customer.Phone == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing phone"})
	}

	if customer.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing password"})
	}

	// Hash the password
	hashedPassword := utils.HashPassword(customer.Password)
	customer.Password = string(hashedPassword)

	// Check if customer already exists
	var existingCustomer models.Customer
	if err := database.DB.Db.Where("phone = ?", customer.Phone).First(&existingCustomer).Error; err == nil {
		return c.Status(400).JSON(fiber.Map{"error": "Customer already exists"})
	}

	// Create the customer in a goroutine
	go func() {
		if err := database.DB.Db.Create(&customer).Error; err != nil {
			// Handle error in goroutine
			// fmt.Println("Error creating customer:", err)
			c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
			return
		}
	}()

	go func() {
		sms.SendSMS(customer.Phone, "Welcome to our e-commerce platform")
	}()

	// return c.Status(200).JSON(customer)
	return c.Status(200).JSON(fiber.Map{"message": "Sign up successful"})
}

func Login(c *fiber.Ctx) error {
	var loginData struct {
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&loginData); err != nil {
		if _, ok := err.(*json.UnmarshalTypeError); ok {
			if strings.Contains(err.Error(), "phone") {
				return c.Status(400).JSON(fiber.Map{"error": "Missing phone number of type string"})
			}
			if strings.Contains(err.Error(), "password") {
				return c.Status(400).JSON(fiber.Map{"error": "Missing password of type string"})
			}
		}
		if strings.Contains(err.Error(), "unexpected end of JSON input") {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON input"})
		}
		if errors.Is(err, fiber.ErrUnprocessableEntity) {
			return c.Status(422).JSON(fiber.Map{"error": "Unprocessable Entity"})
		}
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	if loginData.Phone == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing phone"})
	}

	if loginData.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing password"})
	}

	// Retrieve user from the database
	user, err := GetUserByPhone(loginData.Phone)
	if err != nil {
		return err
	}

	// Check if the user exists and the password matches
	if user == nil || !utils.CheckPasswordHash(loginData.Password, user.Password) {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Execute GetAccessToken function to get the access token
	accessToken, err := auth.GetAccessToken()

	// Extract access token from response
	if accessToken == "" {
		fmt.Println("Access token not found in Hydra token creation response")
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}

	if err != nil {
		fmt.Println("Error getting access token:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}

	// Set the access token in the response headers
	c.Set("Authorization", "Bearer "+accessToken)

	// Return a success message
	return c.Status(400).JSON(fiber.Map{
		"message":      "Login successful",
		"access_token": accessToken,
		"token_type":   "bearer",
	})
}

// GetCustomer retrieves information of currently authorized user
func GetCustomer(c *fiber.Ctx) error {
	// Retrieve user information from the context
	user := c.Locals("user").(*models.Customer)

	// Return user information
	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
	// customer := new(models.Customer)
	return nil
}

func CreateCart(c *fiber.Ctx) error {
	order := new(models.Order)

	// Set the customer_id from the authorized user
	user := c.Locals("user").(*models.Customer)
	order.CustomerID = user.ID

	// Error check fields
	if err := c.BodyParser(order); err != nil {
		if _, ok := err.(*json.UnmarshalTypeError); ok {
			if strings.Contains(err.Error(), "product_id") {
				return c.Status(400).JSON(fiber.Map{"error": "Missing product_id of type int"})
			}
			if strings.Contains(err.Error(), "quantity") {
				return c.Status(400).JSON(fiber.Map{"error": "Missing quantity of type int"})
			}
		}
		// return c.Status(400).SendString(err.Error())
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	if order.ProductID <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Missing product_id"})
	}

	if order.Quantity <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Missing quantity"})
	}

	// Retrieve product from the database
	var product models.Product
	if err := database.DB.Db.Where("id = ?", order.ProductID).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(400).JSON(fiber.Map{"error": "Product not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}

	// check if product is available
	if product.Stock < order.Quantity {
		return c.Status(400).JSON(fiber.Map{"error": "Product not available"})
	}

	// Set the product_id and amount
	order.ProductID = product.ID
	order.Amount = order.Quantity * product.Price

	// Set the time and status
	order.Status = "cart"

	// Create the order in a goroutine
	go func() {
		if err := database.DB.Db.Create(&order).Error; err != nil {
			// Handle error in goroutine
			// fmt.Println("Error creating order:", err)
			c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
			return
		}
	}()

	return c.Status(200).JSON(order)
}

func GetCart(c *fiber.Ctx) error {
	// Retrieve user information from the context
	user := c.Locals("user").(*models.Customer)

	// Retrieve cart items from the database
	var cartItems []models.Cart
	if err := database.DB.Db.Table("orders").Select("product_id, quantity").Where("customer_id = ? AND status = ?", user.ID, "cart").Scan(&cartItems).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}

	// Return cart items
	return c.JSON(cartItems)
}

func UpdateCart(c *fiber.Ctx) error {
	// Retrieve user information from the context
	user := c.Locals("user").(*models.Customer)

	// Retrieve product_id and quantity from the request
	var cartItem models.Cart
	if err := c.BodyParser(&cartItem); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Retrieve the order from the database
	var order models.Order
	if err := database.DB.Db.Where("id = ? AND customer_id = ? AND status = ?", c.Params("id"), user.ID, "cart").First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(400).JSON(fiber.Map{"error": "Order not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}

	// Update the order in the database
	if err := database.DB.Db.Model(&order).Updates(models.Order{Quantity: cartItem.Quantity}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.Status(200).JSON(order)
}

func DeleteCart(c *fiber.Ctx) error {
	// Retrieve user information from the context
	user := c.Locals("user").(*models.Customer)

	// Retrieve the order from the database
	var order models.Order
	if err := database.DB.Db.Where("id = ? AND customer_id = ? AND status = ?", c.Params("id"), user.ID, "cart").First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(400).JSON(fiber.Map{"error": "Order not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}

	// Delete the order from the database
	if err := database.DB.Db.Delete(&order).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Order deleted"})
}

func CreateOrder(c *fiber.Ctx) error {
	// Retrieve user information from the context
	user := c.Locals("user").(*models.Customer)

	order := new(models.Order)

	// Error check fields
	if err := c.BodyParser(order); err != nil {
		if _, ok := err.(*json.UnmarshalTypeError); ok {
			if strings.Contains(err.Error(), "product_id") {
				return c.Status(400).JSON(fiber.Map{"error": "Missing product_id of type int"})
			}
		}
		// return c.Status(400).SendString(err.Error())
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	if order.ProductID <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Missing product_id"})
	}

	// Retrieve product from the database
	var product models.Product
	if err := database.DB.Db.Where("id = ?", order.ProductID).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(400).JSON(fiber.Map{"error": "Product not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}

	// check if product is available
	if product.Stock < order.Quantity {
		return c.Status(400).JSON(fiber.Map{"error": "Product not available"})
	}

	// Set the product_id and amount
	order.ProductID = product.ID
	order.Amount = order.Quantity * product.Price

	// Set the time and status
	order.Time = time.Now().Format("2006-01-02 15:04:05")
	order.Status = "ordered"

	// Create the order in a goroutine
	go func() {
		if err := database.DB.Db.Create(&order).Error; err != nil {
			// Handle error in goroutine
			fmt.Println("Error creating order:", err)
			c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
			return
		}
	}()

	go func() {
		sms.SendSMS(user.Phone, "Order successful")
	}()

	// reduce the stock of the product
	product.Stock -= order.Quantity
	if err := database.DB.Db.Model(&product).Updates(models.Product{Stock: product.Stock}).Error; err != nil {
		fmt.Println("Error updating product stock:", err)
	}

	return c.Status(200).JSON(order)
}
