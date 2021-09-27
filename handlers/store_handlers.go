package handlers

import (
	"encoding/json"
	"github.com/rahularcota/DIY1/dal/models"
	"github.com/rahularcota/DIY1/dal/repos"
	"github.com/rahularcota/DIY1/db_util"
	"net/http"
)

func CreateStore(w http.ResponseWriter, r *http.Request) {
	var s models.Store
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&s); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	repo := repos.GetStoreRepo()
	if err := repo.CreateStore(db_util.GetDB(), &s); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, s)
}
