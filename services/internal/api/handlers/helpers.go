package handlers

import (
	"errors"

	"github.com/leroysb/go_kubernetes/internal/database"
	"github.com/leroysb/go_kubernetes/internal/database/models"
	"gorm.io/gorm"
)

func GetUserByPhone(phone string) (*models.Customer, error) {
	var user models.Customer
	if err := database.DB.Db.Where("phone = ?", phone).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Return nil if user not found
			return nil, nil
		}
		// Return error for other database errors
		return nil, err
	}
	return &user, nil
}
