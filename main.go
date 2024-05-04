package main

import (
	"log"
)

func main() {
	db, err := NewPostgreSQLStorage()
	if err != nil {
		log.Fatal(err)
	}

	apiServer := NewAPIServer(":8080", db)
	apiServer.Run()

}
