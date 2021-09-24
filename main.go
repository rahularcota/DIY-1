package main

import (
	"github.com/rahularcota/DIY1/db_util"
	"log"
	"net/http"
)

func main() {
	db_util.InitializeDB()
	db_util.MigrateDB()
	router := InitializeRouter()

	//a.Run(":8010")
	log.Fatal(http.ListenAndServe(":8010", router))
}