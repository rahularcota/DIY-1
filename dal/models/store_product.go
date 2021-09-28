package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type StoreProduct struct {
	gorm.Model
	StoreId     int  `json:"store_id"`
	ProductId   int  `json:"product_id"`
	IsAvailable bool `json:"is_available"`
}