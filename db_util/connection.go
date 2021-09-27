package db_util

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/rahularcota/DIY1/dal/models"
	"os"
)

var dbCon *gorm.DB

func InitializeDB(){
	host := os.Getenv("HOST")
	dbPort := os.Getenv("DBPORT")
	user := os.Getenv("USER")
	dbname := os.Getenv("NAME")
	dbpassword := os.Getenv("PASSWORD")

	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", host, user, dbname, dbpassword, dbPort)

	db, err := gorm.Open("postgres", dbURI)

	if err != nil {
		panic(err)
	} else {
		fmt.Println("Connected to database successfully")
	}
	dbCon = db
}

func MigrateDB() {
	dbCon.AutoMigrate(&models.Product{})
	dbCon.AutoMigrate(&models.Store{})
	dbCon.AutoMigrate(&models.StoreProduct{})
}

func GetDB() *gorm.DB{
	return dbCon
}
