package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Product struct {
	gorm.Model

	Name  string  `json:"name"`
	Price float64 `json:"price"`
}