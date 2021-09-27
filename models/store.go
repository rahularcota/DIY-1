package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Store struct {
	gorm.Model
	Name string `json:"name"`
}

func(s *Store) CreateStore(db *gorm.DB) error {
	txn := db.Create(&s)

	return txn.Error
}

func(s *Store) GetStore(db *gorm.DB) (bool, error){

	txn := db.First(&s)
	if txn.Error != nil {
		return false, txn.Error
	}

	return true, nil
}
