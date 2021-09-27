package test

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/rahularcota/DIY1/models"
	"net/http"
	"net/http/httptest"
	"strconv"
)

func clearTable(db *gorm.DB) {
	db.Unscoped().Delete(&models.Product{})
	db.Unscoped().Delete(&models.Store{})
	db.Unscoped().Delete(&models.StoreProduct{})
	db.Exec("ALTER SEQUENCE products_id_seq RESTART WITH 1")
	db.Exec("ALTER SEQUENCE stores_id_seq RESTART WITH 1")
}

func executeRequest(req *http.Request, router *mux.Router) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr
}

func addProducts(count int, db *gorm.DB) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		db.Create(&models.Product{Name: strconv.Itoa(i), Price: float64(10.0 * i)})
	}
}

func addStores(count int, db *gorm.DB) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		//db.Exec("INSERT INTO stores(name) VALUES($1)", "Store "+strconv.Itoa(i))
		db.Create(&models.Store{Name: strconv.Itoa(i)})
	}
}

func addStoreProducts(storeId int, productsCount int, db *gorm.DB) {
	if productsCount<1 {
		productsCount = 1
	}

	for i := 0; i < productsCount; i++ {
		db.Create(&models.Product{Name: strconv.Itoa(i), Price: float64(10.0 * i)})
		db.Create(&models.StoreProduct{StoreId: storeId, ProductId: i+1, IsAvailable: true})
		//db.Exec("INSERT INTO storeProduct(store_id, product_id, is_available) VALUES ($1, $2, true)", storeId, i+1)
	}
}
