package test

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/rahularcota/DIY1/db_util"
	"github.com/rahularcota/DIY1/router"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type ProductTestSuite struct {
	suite.Suite
	rtr *mux.Router
	db *gorm.DB
}

func TestProductTestSuite(t *testing.T) {
	suite.Run(t, new(ProductTestSuite))
}

func (suite *ProductTestSuite) SetupSuite() {
	db_util.InitializeDB()
	db_util.MigrateDB()
	suite.rtr = router.InitializeRouter()
	suite.db = db_util.GetDB()
}

func (suite *ProductTestSuite) BeforeTest(suiteName, testName string) {
	clearTable(suite.db)
}

func (suite *ProductTestSuite) TestCreateProduct() {
	var jsonStr = []byte(`{"name":"test product", "price": 11.22}`)
	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req, suite.rtr)
	//checkResponseCode(t, http.StatusCreated, response.Code)
	suite.Equal(http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	suite.Equal(m["name"], "test product")
	suite.Equal(m["price"], 11.22)
	suite.Equal(m["ID"], 1.0)
}

func (suite *ProductTestSuite) TestGetNonExistentProduct() {

	req, _ := http.NewRequest("GET", "/product/11", nil)
	response := executeRequest(req, suite.rtr)
	suite.Equal(http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	suite.Equal(m["error"], "Product not found")
}

func (suite *ProductTestSuite) TestGetProduct() {
	addProducts(1, suite.db)

	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := executeRequest(req, suite.rtr)

	suite.Equal(http.StatusOK, response.Code)
}

func (suite *ProductTestSuite) TestUpdateProduct() {
	addProducts(1, suite.db)

	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := executeRequest(req, suite.rtr)
	var originalProduct map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalProduct)

	var jsonStr = []byte(`{"name":"test product - updated name", "price": 11.22}`)
	req, _ = http.NewRequest("PUT", "/product/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response = executeRequest(req, suite.rtr)

	suite.Equal(http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	suite.Equal(m["name"], "test product - updated name")
	suite.Equal(m["price"], 11.22)
	suite.Equal(m["ID"], originalProduct["ID"])
}

func (suite *ProductTestSuite) TestDeleteProduct() {
	addProducts(1, suite.db)

	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := executeRequest(req, suite.rtr)
	suite.Equal(http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/product/1", nil)
	response = executeRequest(req, suite.rtr)

	suite.Equal(http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/product/1", nil)
	response = executeRequest(req, suite.rtr)
	suite.Equal(http.StatusNotFound, response.Code)
}
