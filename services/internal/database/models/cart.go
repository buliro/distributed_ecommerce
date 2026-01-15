package models

type Cart struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}
