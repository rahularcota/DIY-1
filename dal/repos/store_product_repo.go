package repos

import (
	"github.com/jinzhu/gorm"
	"github.com/rahularcota/DIY1/dal/models"
)

var spRepo StoreProductRepo

type IStoreProductRepo interface {
	AddStoreProducts(db *gorm.DB, products *[] models.Product, storeId int) ([]models.Product, error)
	GetStoreProducts(db *gorm.DB, id int) (*models.Product, error)
}

type StoreProductRepo struct {}

func InitStoreProductRepo() {
	spRepo = StoreProductRepo{}
}

func GetStoreProductRepo() *StoreProductRepo {
	return &spRepo
}

func (repo *StoreProductRepo) AddStoreProducts(db *gorm.DB, products *[] models.Product, storeId int) ([]models.Product, error) {
	addedProducts := make([]models.Product, 0)
	productRepo := GetProductRepo()
	for _, p := range *products {
		if err := productRepo.CreateProduct(db, &p); err == nil {
			txn := db.Create(&models.StoreProduct{StoreId: storeId, ProductId: int(p.ID), IsAvailable: true})
			if txn.Error == nil {
				addedProducts = append(addedProducts, p)
			}
		}
	}
	return addedProducts, nil
}

func (repo *StoreProductRepo) GetStoreProducts(db *gorm.DB, id int) ([]models.Product, error) {
	var products []models.Product
	txn := db.Raw("SELECT products.id, products.name, products.price FROM products inner join store_products on products.id = store_products.product_id WHERE store_products.store_id = ? AND store_products.is_available = true", id).Scan(&products)

	return products, txn.Error
}
