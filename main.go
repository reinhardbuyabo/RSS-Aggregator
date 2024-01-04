package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/reinhardbuyabo/RSS-Aggregator/internal/database"
)

// 1. apiConfig that we can pass to our handler and then connect to our database
type apiConfig struct {
	DB *database.Queries
}

func main() {
	fmt.Println("Hello World")

	// A way to read PORT into main.go ...

	godotenv.Load()

	portString := os.Getenv("PORT")

	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}

	// importing database driver
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the environment")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to database:", err)
	}

	// passing connection to database.New();
	queries := database.New(conn)

	apiCfg := apiConfig{
		DB: queries,
	}

	router := chi.NewRouter()

	// Cross Origin Resource Sharing
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTION"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Respond if the server is alive and running

	// v1Router := chi.NewRouter();
	// v1Router.HandleFunc()

	v1Router := chi.NewRouter()                // to mount on line 42 ... to slash v1 ... /v1
	v1Router.Get("/healthz", handlerReadiness) // connecting path to the function in handler_readiness ... // full path /v1/healthz
	v1Router.Get("/err", handlerErr)
	v1Router.Post("/users", apiCfg.handlerCreateUser)

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server starting on port %v", portString)
	err = srv.ListenAndServe() // returns an error ... Nothing should ever be returned from server.

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Port: ", portString)
}
