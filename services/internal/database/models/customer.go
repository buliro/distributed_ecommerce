package models

import "gorm.io/gorm"

type Customer struct {
	gorm.Model
	Name     string `json:"name" gorm:"text;not null;default:null"`
	Phone    string `json:"phone" gorm:"text;not null;unique"`
	Password string `json:"password" gorm:"text;not null;default:null"`
}
