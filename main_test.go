package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/rahularcota/DIY1/db_util"
	"github.com/rahularcota/DIY1/models"
	"os"
	"testing"

	"bytes"
	"encoding/json"
	// tom: for TestEmptyTable and next functions (no go get is required"
	"net/http"
	// "net/url"
	"net/http/httptest"
	"strconv"
	// "io/ioutil"

)

var router *mux.Router
var db *gorm.DB

func TestMain(m *testing.M) {
	db_util.InitializeDB()
	db_util.MigrateDB()
	router = InitializeRouter()
	db = db_util.GetDB()

	code := m.Run()

	clearTable()

	os.Exit(code)
}

func clearTable() {
	db.Unscoped().Delete(&models.Product{})
	db.Unscoped().Delete(&models.Store{})
	db.Unscoped().Delete(&models.StoreProduct{})
	db.Exec("ALTER SEQUENCE products_id_seq RESTART WITH 1")
	db.Exec("ALTER SEQUENCE stores_id_seq RESTART WITH 1")
}


func TestEmptyTable(t *testing.T) {
	//fmt.Println("before")
	clearTable()
	//fmt.Println("After")

	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		return
	}
	t.Fatal(fmt.Sprintf("%v != %v", a, b))
}

func addProducts(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		//db.Exec("INSERT INTO products(name, price) VALUES($1, $2)", "Product "+strconv.Itoa(i), (i+1.0)*10)
		db.Create(&models.Product{Name: strconv.Itoa(i), Price: float64(10.0 * i)})
	}
}

func addStores(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		//db.Exec("INSERT INTO stores(name) VALUES($1)", "Store "+strconv.Itoa(i))
		db.Create(&models.Store{Name: strconv.Itoa(i)})
	}
}

func addStoreProducts(storeId int, productsCount int) {
	if productsCount<1 {
		productsCount = 1
	}

	for i := 0; i < productsCount; i++ {
		db.Create(&models.Product{Name: strconv.Itoa(i), Price: float64(10.0 * i)})
		db.Create(&models.StoreProduct{StoreId: storeId, ProductId: i+1, IsAvailable: true})
		//db.Exec("INSERT INTO storeProduct(store_id, product_id, is_available) VALUES ($1, $2, true)", storeId, i+1)
	}
}


// product tests
func TestCreateProduct(t *testing.T) {
	clearTable()

	var jsonStr = []byte(`{"name":"test product", "price": 11.22}`)
	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	fmt.Println(m)

	if m["name"] != "test product" {
		t.Errorf("Expected product name to be 'test product'. Got '%v'", m["name"])
	}

	if m["price"] != 11.22 {
		t.Errorf("Expected product price to be '11.22'. Got '%v'", m["price"])
	}

	// the id is compared to 1.0 because JSON unmarshaling converts numbers to
	// floats, when the target is a map[string]interface{}
	if m["ID"] != 1.0 {
		t.Errorf("Expected product ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetNonExistentProduct(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/product/11", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Product not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Product not found'. Got '%s'", m["error"])
	}
}

func TestGetProduct(t *testing.T) {
	clearTable()
	addProducts(1)

	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdateProduct(t *testing.T) {

	clearTable()
	addProducts(1)

	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := executeRequest(req)
	var originalProduct map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalProduct)

	var jsonStr = []byte(`{"name":"test product - updated name", "price": 11.22}`)
	req, _ = http.NewRequest("PUT", "/product/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	// req, _ = http.NewRequest("PUT", "/product/1", bytes.NewBuffer(payload))
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalProduct["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalProduct["id"], m["id"])
	}

	if m["name"] == originalProduct["name"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalProduct["name"], m["name"], m["name"])
	}

	if m["price"] == originalProduct["price"] {
		t.Errorf("Expected the price to change from '%v' to '%v'. Got '%v'", originalProduct["price"], m["price"], m["price"])
	}
}

func TestDeleteProduct(t *testing.T) {
	clearTable()
	addProducts(1)

	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/product/1", nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/product/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}


//store tests
func TestCreateStore(t *testing.T) {
	clearTable()

	var jsonStr = []byte(`{"name":"test store"}`)
	req, _ := http.NewRequest("POST", "/store", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "test store" {
		t.Errorf("Expected store name to be 'test store'. Got '%v'", m["name"])
	}

	if m["ID"] != 1.0 {
		t.Errorf("Expected store ID to be '1'. Got '%v'", m["id"])
	}
}

func TestCheckStoreExistence(t *testing.T) {
	clearTable()
	addStores(1)
	s := models.Store{}
	s.ID = 1
	exists, _ := s.CheckStoreExistence(db)
	assertEqual(t, true, exists)
}

func TestCheckStoreNonExistence(t *testing.T) {
	clearTable()
	s := models.Store{}
	s.ID = 1
	exists, _ := s.CheckStoreExistence(db)
	assertEqual(t, false, exists)
}


//store-product tests

func TestGetStoreProducts(t *testing.T) {
	clearTable()
	addStores(1)
	addStoreProducts(1, 1)

	req, _ := http.NewRequest("GET", "/store/1/products", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestGetStoreProductsInvalidID(t *testing.T) {
	clearTable()
	addStoreProducts(1, 1)

	req, _ := http.NewRequest("GET", "/store/a_b/products", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestAddStoreProducts(t *testing.T) {
	clearTable()
	addStores(1)
	var jsonStr = []byte(`[{"name":"test product", "price": 11.22}]`)
	req, _ := http.NewRequest("POST", "/store/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	resp := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, resp.Code)

	products := make([]models.Product,0)
	json.Unmarshal(resp.Body.Bytes(), &products)

	if len(products)!=1 {
		t.Errorf("Product not added to store")
	}

	if products[0].Name != "test product" {
		t.Errorf("Expected product name to be 'test product'. Got '%v'", products[0].Name)
	}

	if products[0].Price != 11.22 {
		t.Errorf("Expected product price is 11.22. Got '%v'", products[0].Price)
	}

	if products[0].ID != 1 {
		t.Errorf("Expected product ID to be 1. Got '%v'", products[0].ID)
	}
}

func TestAddStoreProductInvalidStoreID(t *testing.T) {
	clearTable()
	var jsonStr = []byte(`[{"name":"test product", "price": 11.22}]`)
	req, _ := http.NewRequest("POST", "/store/a_b", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	resp := executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, resp.Code)
}

func TestAddStoreProductNonExistentStoreID(t *testing.T) {
	clearTable()
	var jsonStr = []byte(`[{"name":"test product", "price": 11.22}]`)
	req, _ := http.NewRequest("POST", "/store/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	resp := executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, resp.Code)
}

func TestAddStoreProductInvalidPayload(t *testing.T) {
	clearTable()
	addStores(1)
	var jsonStr = []byte(`[{"name":123, "price": 11.22}]`)
	req, _ := http.NewRequest("POST", "/store/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	resp := executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, resp.Code)
}

func TestGetStoreProductsNonExistentStoreID(t *testing.T) {
	clearTable()
	req, _ := http.NewRequest("GET", "/store/1/products", nil)

	resp := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, resp.Code)

}