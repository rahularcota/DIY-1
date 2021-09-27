package test

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/rahularcota/DIY1/dal/models"
	"github.com/rahularcota/DIY1/dal/repos"
	"github.com/rahularcota/DIY1/db_util"
	"github.com/rahularcota/DIY1/router"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type StoreProductTestSuite struct {
	suite.Suite
	rtr *mux.Router
	db *gorm.DB
}

func TestStoreProductTestSuite(t *testing.T) {
	suite.Run(t, new(StoreProductTestSuite))
}

func (suite *StoreProductTestSuite) SetupSuite() {
	db_util.InitializeDB()
	db_util.MigrateDB()
	repos.InitProductRepo()
	repos.InitStoreRepo()
	repos.InitStoreProductRepo()
	suite.rtr = router.InitializeRouter()
	suite.db = db_util.GetDB()
}

func (suite *StoreProductTestSuite) BeforeTest(suiteName, testName string) {
	clearTable(suite.db)
}

func (suite *StoreProductTestSuite) TestGetStoreProducts() {
	addStores(1, suite.db)
	addStoreProducts(1, 1, suite.db)

	req, _ := http.NewRequest("GET", "/store/1/products", nil)
	response := executeRequest(req, suite.rtr)

	suite.Equal(http.StatusOK, response.Code)
}

func (suite *StoreProductTestSuite) TestGetStoreProductsInvalidID() {
	addStoreProducts(1, 1, suite.db)

	req, _ := http.NewRequest("GET", "/store/a_b/products", nil)
	response := executeRequest(req, suite.rtr)

	suite.Equal(http.StatusBadRequest, response.Code)
}

func (suite *StoreProductTestSuite) TestAddStoreProducts() {
	addStores(1, suite.db)
	var jsonStr = []byte(`[{"name":"test product", "price": 11.22}]`)
	req, _ := http.NewRequest("POST", "/store/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req, suite.rtr)
	suite.Equal(http.StatusCreated, response.Code)

	products := make([]models.Product,0)
	json.Unmarshal(response.Body.Bytes(), &products)

	suite.Equal(len(products), 1)
	suite.Equal(products[0].Name, "test product")
	suite.Equal(products[0].Price, 11.22)
	suite.Equal(int(products[0].ID), 1)
}

func (suite *StoreProductTestSuite) TestAddStoreProductInvalidStoreID() {
	var jsonStr = []byte(`[{"name":"test product", "price": 11.22}]`)
	req, _ := http.NewRequest("POST", "/store/a_b", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req, suite.rtr)
	suite.Equal(http.StatusBadRequest, response.Code)
}

func (suite *StoreProductTestSuite) TestAddStoreProductNonExistentStoreID() {
	var jsonStr = []byte(`[{"name":"test product", "price": 11.22}]`)
	req, _ := http.NewRequest("POST", "/store/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req, suite.rtr)
	suite.Equal(http.StatusNotFound, response.Code)
}

func (suite *StoreProductTestSuite) TestAddStoreProductInvalidPayload() {
	addStores(1, suite.db)
	var jsonStr = []byte(`[{"name":123, "price": 11.22}]`)
	req, _ := http.NewRequest("POST", "/store/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req, suite.rtr)
	suite.Equal(http.StatusBadRequest, response.Code)
}

func (suite *StoreProductTestSuite) TestGetStoreProductsNonExistentStoreID() {
	req, _ := http.NewRequest("GET", "/store/1/products", nil)

	response := executeRequest(req, suite.rtr)

	suite.Equal(http.StatusNotFound, response.Code)

}