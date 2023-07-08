package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("hello world")

	godotenv.Load(".env") // takes the environment variable from .env file and brings it here:

	portString := os.Getenv("PORT") // exported fn Getenv from os library in the stdlib.

	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}

	router := chi.NewRouter()

	srv := &http.Server{ // pointer to http server
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server starting on port %v", portString)
	err := srv.ListenAndServe() // returns an error, we capture the error here... ListenAndServe will block ... Our code stops here and starts handling http requests ... Nothing should be returned from the server, since it will run forever
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Port: ", portString)
}
