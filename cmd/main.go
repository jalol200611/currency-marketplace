package main

import (
	"currency/database"
	"currency/routes"
	"log"
	"net/http"
)

func main() {
	router := routes.SetupRoutes()

	log.Println("🚀Сервер запущен на :8080")

	database.Connect()

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}
