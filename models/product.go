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

func (p *Product) GetProduct(db *gorm.DB) error {
	txn := db.First(&p, p.ID)
	return txn.Error
	//return db.QueryRow("SELECT name, price FROM products WHERE id=$1",
	//	p.ID).Scan(&p.Name, &p.Price)
}

func (p *Product) UpdateProduct(db *gorm.DB) error {
	//_, err :=
	//	db.Exec("UPDATE products SET name=$1, price=$2 WHERE id=$3",
	//		p.Name, p.Price, p.ID)
	//
	//return err
	txn := db.Save(&p)

	return txn.Error
}

func (p *Product) DeleteProduct(db *gorm.DB) error {
	//_, err := db.Exec("DELETE FROM products WHERE id=$1", p.ID)
	//
	//return err
	txn := db.Delete(&p, p.ID)
	return txn.Error
}

func (p *Product) CreateProduct(db *gorm.DB) error {
	//err := db.QueryRow(
	//	"INSERT INTO products(name, price) VALUES($1, $2) RETURNING id",
	//	p.Name, p.Price).Scan(&p.ID)
	//
	//if err != nil {
	//	return err
	//}
	txn := db.Create(&p)

	return txn.Error
}