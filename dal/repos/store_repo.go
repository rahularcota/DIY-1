package repos

import (
	"github.com/jinzhu/gorm"
	"github.com/rahularcota/DIY1/dal/models"
)

var storeRepo StoreRepo

type IStoreRepo interface {
	CreateStore(db *gorm.DB, store *models.Store) error
	GetStore(db *gorm.DB, id int) (bool, error)
}

type StoreRepo struct {}

func InitStoreRepo() {
	storeRepo = StoreRepo{}
}

func GetStoreRepo() *StoreRepo {
	return &storeRepo
}

func (repo *StoreRepo) CreateStore(db *gorm.DB, s *models.Store) error {
	txn := db.Create(&s)
	return txn.Error
}

func (repo *StoreRepo) GetStore(db *gorm.DB, id int) (bool, error) {
	s := &models.Store{}
	txn := db.First(&s, id)
	if txn.Error != nil {
		return false, txn.Error
	}
	return true, nil
}
