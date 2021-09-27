package main

import (
	"github.com/rahularcota/DIY1/db_util"
	"github.com/rahularcota/DIY1/router"
	"log"
	"net/http"
)

func main() {
	db_util.InitializeDB()
	db_util.MigrateDB()
	router := router.InitializeRouter()

	//a.Run(":8010")
	log.Fatal(http.ListenAndServe(":8010", router))
}