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

//func (p *Product) GetProduct(db *gorm.DB) error {
//	txn := db.First(&p, p.ID)
//	return txn.Error
//}
//
//func (p *Product) UpdateProduct(db *gorm.DB) error {
//	txn := db.Save(&p)
//	return txn.Error
//}
//
//func (p *Product) DeleteProduct(db *gorm.DB) error {
//	txn := db.Delete(&p, p.ID)
//	return txn.Error
//}
//
//func (p *Product) CreateProduct(db *gorm.DB) error {
//	txn := db.Create(&p)
//	return txn.Error
//}