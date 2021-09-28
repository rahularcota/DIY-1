package repos

import (
	"github.com/jinzhu/gorm"
	"github.com/rahularcota/DIY1/dal/models"
)

var productRepo ProductRepo

type IProductRepo interface {
	CreateProduct(db *gorm.DB, product *models.Product) error
	UpdateProduct(db *gorm.DB, product *models.Product) (*models.Product, error)
	DeleteProduct(db *gorm.DB, id int) error
	GetProduct(db *gorm.DB, id int) (*models.Product, error)
}

type ProductRepo struct {

}

func InitProductRepo() {
	productRepo = ProductRepo{}
}

func GetProductRepo() *ProductRepo {
	return &productRepo
}

func (repo *ProductRepo) GetProduct(db *gorm.DB, id int) (*models.Product, error) {
	p := &models.Product{}
	txn := db.First(&p, id)
	return p, txn.Error
}

func (repo *ProductRepo) UpdateProduct(db *gorm.DB, p *models.Product) error {
	txn := db.Save(&p)
	return txn.Error
}

func (repo *ProductRepo) DeleteProduct(db *gorm.DB, id int) error {
	txn := db.Delete(&models.Product{}, id)
	return txn.Error
}

func (repo *ProductRepo) CreateProduct(db *gorm.DB, p *models.Product) error {
	txn := db.Create(&p)
	return txn.Error
}
