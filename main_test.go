package main

import (
	"fmt"
	"os"
	"testing"

	// tom: for ensureTableExists
	"log"

	"bytes"
	"encoding/json"
	// tom: for TestEmptyTable and next functions (no go get is required"
	"net/http"
	// "net/url"
	"net/http/httptest"
	"strconv"
	// "io/ioutil"

)

var a App

func TestMain(m *testing.M) {
	a = App{}
	a.Initialize(
		os.Getenv("TEST_DB_USERNAME"),
		os.Getenv("TEST_DB_PASSWORD"),
		os.Getenv("TEST_DB_NAME"))

	ensureTableExists()

	code := m.Run()

	clearTable()

	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(productsTableCreationQuery); err != nil {
		log.Fatal(err)
	}

	if _, err := a.DB.Exec(storeProductTableCreationQuery); err != nil {
		log.Fatal(err)
	}

	if _, err := a.DB.Exec(storeTableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM products")
	a.DB.Exec("ALTER SEQUENCE products_id_seq RESTART WITH 1")

	a.DB.Exec("DELETE FROM storeProduct")

	a.DB.Exec("DELETE FROM stores")
	a.DB.Exec("ALTER SEQUENCE stores_id_seq RESTART WITH 1")
}

const productsTableCreationQuery = `CREATE TABLE IF NOT EXISTS products
(
    id SERIAL,
    name TEXT NOT NULL,
    price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
    CONSTRAINT products_pkey PRIMARY KEY (id)
)`
const storeProductTableCreationQuery = `CREATE TABLE IF NOT EXISTS storeProduct
(
	store_id NUMERIC,
	product_id NUMERIC,
	is_available BOOLEAN
)`
const storeTableCreationQuery = `CREATE TABLE IF NOT EXISTS stores
(
	id SERIAL,
	name TEXT NOT NULL,
	CONSTRAINT stores_pkey PRIMARY KEY (id)
)`


// tom: next functions added later, these require more modules: net/http net/http/httptest
func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/products", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

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
		a.DB.Exec("INSERT INTO products(name, price) VALUES($1, $2)", "Product "+strconv.Itoa(i), (i+1.0)*10)
	}
}

func addStores(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		a.DB.Exec("INSERT INTO stores(name) VALUES($1)", "Store "+strconv.Itoa(i))
	}
}

func addStoreProducts(storeId int, productsCount int) {
	if productsCount<1 {
		productsCount = 1
	}

	for i := 0; i < productsCount; i++ {
		a.DB.Exec("INSERT INTO products(name, price) VALUES($1, $2)", "Product "+strconv.Itoa(i), (i+1.0)*10)
		a.DB.Exec("INSERT INTO storeProduct(store_id, product_id, is_available) VALUES ($1, $2, true)", storeId, i+1)
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

	if m["name"] != "test product" {
		t.Errorf("Expected product name to be 'test product'. Got '%v'", m["name"])
	}

	if m["price"] != 11.22 {
		t.Errorf("Expected product price to be '11.22'. Got '%v'", m["price"])
	}

	// the id is compared to 1.0 because JSON unmarshaling converts numbers to
	// floats, when the target is a map[string]interface{}
	if m["id"] != 1.0 {
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

	if m["id"] != 1.0 {
		t.Errorf("Expected store ID to be '1'. Got '%v'", m["id"])
	}
}

func TestCheckStoreExistence(t *testing.T) {
	clearTable()
	addStores(1)
	s := store{ID: 1}
	exists, _ := s.checkStoreExistence(a.DB)
	assertEqual(t, true, exists)
}

func TestCheckStoreNonExistence(t *testing.T) {
	clearTable()
	s := store{ID: 1}
	exists, _ := s.checkStoreExistence(a.DB)
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

	products := make([]product,0)
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