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

//func (sp *StoreProduct) GetStoreProducts(db *gorm.DB) ([]Product, error) {
//	var products []Product
//	txn := db.Raw("SELECT products.id, products.name, products.price FROM products inner join store_products on products.id = store_products.product_id WHERE store_products.store_id = ? AND store_products.is_available = true", sp.StoreId).Scan(&products)
//
//	return products, txn.Error
//}
//
//func (sp *StoreProduct) AddStoreProducts(db *gorm.DB, products []Product) ([]Product, error) {
//	addedProducts := make([]Product, 0)
//	for _, p := range products {
//		if err := p.CreateProduct(db); err == nil {
//			txn := db.Create(&StoreProduct{StoreId: sp.StoreId, ProductId: int(p.ID), IsAvailable: true})
//			if txn.Error == nil {
//				addedProducts = append(addedProducts, p)
//			}
//		}
//	}
//
//	return addedProducts, nil
//}
