package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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
	feed, err := urlToFeed("https://wagslane.dev/index.xml")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(feed)

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
	db := database.New(conn)

	apiCfg := apiConfig{
		DB: db,
	}

	go startScraping(db, 10, time.Minute) // calling it on a new go routine, so that it doesn't interrupt this main flow, NOTE: startScraping is a long running function[it has an infinite for loop on line 26]

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
	// 1. Call the middlewareAuth function first 2. Get User By Api Key(middleware_auth.go) 3. Get API Key (auth.go) which returns apiKey of user(middleware_auth.go) 4. We Query the database using the gotten apiKey (GetUserByAPIKey()) and pass the resulted row, gotten user to our handler function in handler_user.go (middleware_auth.go)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handleGetUser)) // cannot use apiCfg.handleGetUser (value of type func(w http.ResponseWriter, r *http.Request, user database.User)) as http.HandleFunc value in argument to v1Router.Get

	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handleCreateFeed))
	v1Router.Get("/feeds", apiCfg.handlerGetFeeds)

	v1Router.Get("/posts", apiCfg.middlewareAuth(apiCfg.handleGetPostsForUser))

	v1Router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
	v1Router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollows))
	// we want to dynamically pass the feed follow ID as a parameter to the URL
	v1Router.Delete("/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow))

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
