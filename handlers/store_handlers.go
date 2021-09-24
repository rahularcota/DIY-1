package handlers

import (
	"encoding/json"
	"github.com/rahularcota/DIY1/db_util"
	"github.com/rahularcota/DIY1/models"
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

	if err := s.CreateStore(db_util.GetDB()); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, s)
}
