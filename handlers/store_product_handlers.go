package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rahularcota/DIY1/db_util"
	"github.com/rahularcota/DIY1/models"
	"net/http"
	"strconv"
)

func GetStoreProducts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Store ID")
		return
	}

	s := models.Store{}
	s.ID = uint(id)
	exists, err:= s.GetStore(db_util.GetDB())
	if err != nil  && !exists{
		respondWithError(w, http.StatusNotFound, "Store does not exist")
		return
	}

	sp := models.StoreProduct{StoreId: id}
	products, err := sp.GetStoreProducts(db_util.GetDB())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, products)
}

func AddStoreProducts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Store ID")
		return
	}

	s := models.Store{}
	s.ID = uint(id)
	exists, err:= s.GetStore(db_util.GetDB())
	if err != nil  && !exists{
		respondWithError(w, http.StatusNotFound, "Store does not exist")
		return
	}

	sp := models.StoreProduct{}
	sp.StoreId = id
	sp.IsAvailable = true
	products := make([]models.Product,0)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&products); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if products, err = sp.AddStoreProducts(db_util.GetDB(), products); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, products)
}