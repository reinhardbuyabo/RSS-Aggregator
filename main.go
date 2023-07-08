package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
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

	// Adding cors configuration to our browser
	router.Use(cors.Handler(cors.Options{ // This configuration is telling our server to send a bunch of extra HTTP headers in our responses that will tell browsers that we will you to do the specified below. We're making it so permissive. The server is JSON Rest API
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	v1Router.HandleFunc("/healthz", handlerReadiness) // full path ... /v1/ready ... to cater for breaking changes for our REST API

	router.Mount("/v1", v1Router)

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
