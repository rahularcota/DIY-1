package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rahularcota/DIY1/dal/models"
	"github.com/rahularcota/DIY1/dal/repos"
	"github.com/rahularcota/DIY1/db_util"
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

	repo := repos.GetStoreRepo()
	exists, err:= repo.GetStore(db_util.GetDB(), id)
	if err != nil  && !exists{
		respondWithError(w, http.StatusNotFound, "Store does not exist")
		return
	}

	sprepo := repos.GetStoreProductRepo()
	products, err := sprepo.GetStoreProducts(db_util.GetDB(), id)
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

	repo := repos.GetStoreRepo()
	exists, err:= repo.GetStore(db_util.GetDB(), id)
	if err != nil  && !exists{
		respondWithError(w, http.StatusNotFound, "Store does not exist")
		return
	}

	sprepo := repos.GetStoreProductRepo()
	products := make([]models.Product,0)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&products); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if products, err = sprepo.AddStoreProducts(db_util.GetDB(), &products, id); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, products)
}