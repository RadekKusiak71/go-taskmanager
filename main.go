package main

import (
	"log"

	"github.com/RadekKusiak71/taskmanager/api"
	"github.com/RadekKusiak71/taskmanager/db"
)

func main() {
	db, err := db.NewPostgreSQLStorage()
	if err != nil {
		log.Fatal(err)
	}

	apiServer := api.NewAPIServer(":8080", db)
	apiServer.Run()

}
