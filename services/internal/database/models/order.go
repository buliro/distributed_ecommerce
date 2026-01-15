package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	Customer   Customer `gorm:"foreignKey:CustomerID"`
	CustomerID uint     `json:"customer_id" gorm:"integer;not null;default:null"`
	Product    Product  `gorm:"foreignKey:ProductID"`
	ProductID  uint     `json:"product_id" gorm:"integer;not null;default:null"`
	Quantity   int      `json:"quantity" gorm:"integer;not null;default:null"`
	Amount     int      `json:"amount" gorm:"float;not null;default:null"`
	Time       string   `json:"time" gorm:"text;not null;default:null"`
	Status     string   `json:"status" gorm:"text;not null;default:null"`
}
