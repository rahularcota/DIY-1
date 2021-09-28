package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Store struct {
	gorm.Model
	Name string `json:"name"`
}
