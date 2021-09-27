package test

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/rahularcota/DIY1/db_util"
	"github.com/rahularcota/DIY1/models"
	"github.com/rahularcota/DIY1/router"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type StoreTestSuite struct {
	suite.Suite
	rtr *mux.Router
	db *gorm.DB
}

func TestStoreTestSuite(t *testing.T) {
	suite.Run(t, new(StoreTestSuite))
}

func (suite *StoreTestSuite) SetupSuite() {
	db_util.InitializeDB()
	db_util.MigrateDB()
	suite.rtr = router.InitializeRouter()
	suite.db = db_util.GetDB()
}

func (suite *StoreTestSuite) BeforeTest(suiteName, testName string) {
	clearTable(suite.db)
}

func (suite *StoreTestSuite) TestCreateStore() {
	var jsonStr = []byte(`{"name":"test store"}`)
	req, _ := http.NewRequest("POST", "/store", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req, suite.rtr)
	suite.Equal(http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	suite.Equal(m["name"], "test store")
	suite.Equal(m["ID"], 1.0)
}

func (suite *StoreTestSuite) TestGetStore() {
	addStores(1, suite.db)
	s := models.Store{}
	s.ID = 1
	exists, _ := s.GetStore(suite.db)
	suite.Equal(true, exists)
}

func (suite *StoreTestSuite) TestGetNonExistentStore() {
	s := models.Store{}
	s.ID = 1
	exists, _ := s.GetStore(suite.db)
	suite.Equal(false, exists)
}
