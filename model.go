// model.go

package main

import (
	"database/sql"
	// tom: errors is removed once functions are implemented
	// "errors"
)


// tom: add backticks to json
type product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// tom: these are added after tdd tests
func (p *product) getProduct(db *sql.DB) error {
	return db.QueryRow("SELECT name, price FROM products WHERE id=$1",
		p.ID).Scan(&p.Name, &p.Price)
}

func (p *product) updateProduct(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE products SET name=$1, price=$2 WHERE id=$3",
			p.Name, p.Price, p.ID)

	return err
}

func (p *product) deleteProduct(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM products WHERE id=$1", p.ID)

	return err
}

func (p *product) createProduct(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO products(name, price) VALUES($1, $2) RETURNING id",
		p.Name, p.Price).Scan(&p.ID)

	if err != nil {
		return err
	}

	return nil
}

func getProducts(db *sql.DB, start, count int) ([]product, error) {
	rows, err := db.Query(
		"SELECT id, name,  price FROM products LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	products := []product{}

	for rows.Next() {
		var p product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}


type store struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}


func(s *store) createStore(db *sql.DB) error {
	err := db.QueryRow("INSERT INTO stores(name) Values($1) RETURNING id", s.Name).Scan(&s.ID)

	if err != nil {
		return err
	}

	return nil
}

func(s *store) checkStoreExistence(db *sql.DB) (bool, error){
	rows, err := db.Query("SELECT * FROM stores WHERE id = $1", s.ID)

	if err != nil {
		return false, err
	}

	for rows.Next() {
		return true, nil
	}

	return false, nil
}


type storeProduct struct {
	StoreId     int  `json:"store_id"`
	ProductId   int  `json:"product_id"`
	IsAvailable bool `json:"is_available"`
}


func (s *storeProduct) getStoreProducts(db *sql.DB) ([]product, error){
	rows, err := db.Query(
		"SELECT products.id, products.name, products.price FROM products inner join storeProduct on products.id = storeProduct.product_id WHERE storeProduct.store_id = $1 AND storeProduct.is_available = true", s.StoreId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	products := []product{}

	for rows.Next() {
		var p product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}


func (s *storeProduct) addStoreProducts(db *sql.DB, products []product) ([]product, error) {
	addedProducts := make([]product, 0)
	for _, p := range products {
		if err := p.createProduct(db); err==nil {
			_, err2 := db.Exec("INSERT INTO storeProduct VALUES ($1, $2, true)", s.StoreId, p.ID)
			if err2 == nil {
				addedProducts = append(addedProducts, p)
			}
		}
	}

	return addedProducts, nil
}
