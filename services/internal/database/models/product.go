package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name  string `json:"name" gorm:"text;not null;default:null"`
	Price int    `json:"price" gorm:"float;not null;default:null"`
	Stock int    `json:"stock" gorm:"integer;not null;default:null"`
}
