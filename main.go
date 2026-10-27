package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Dnts77/rss-aggregator-go/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {

	godotenv.Load(".env") //	External go module that allows to pull .env variables

	portString := os.Getenv("PORT") //	os.Getenv -> gets the value of a variable in .env by its key

	if portString == "" {
		log.Fatal("PORT is not found in the environment") //	Logging the error
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the environment") //	Logging the error
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to database:", err)
	}

	apiCfg := apiConfig{
		DB: database.New(conn),
	}

	router := chi.NewRouter() //	New Router object

	router.Use(cors.Handler(cors.Options{
		//	This configuration tells the server to send a bunch of extra HTTP Headers in our responses
		// 	that will tell browsers that they are allowed to send requests to "https://" and "http://" versions,
		//	they're allowd to use the listed methods,
		//	they're allowed to send any headers, etc
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)                               //	Connecting the "handlerReadiness" function to the "/healthz" path - Only for GET requests
	v1Router.Get("/error", handlerErr)                                       //	Handling errors with the "handleErr" function
	v1Router.Post("/users", apiCfg.handlerCreateUser)                        //	Creating users' route
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))     //	Getting the user by their API key. Works through a middleware to dry up the code
	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed)) //	Creating a feeds' route by the middleware and the handlerCreateFeed method
	v1Router.Get("/feeds", apiCfg.handlerGetFeeds)

	router.Mount("/v1", v1Router) //	Nestint the "v1Router" in the "/v1" path. So the full path to this request will be "/v1/healthz"

	srv := &http.Server{ //	Creating the server -> The Handler is the router itself and the address is a ":" + the portString
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server starting on PORT: %v", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("PORT: ", portString)
}
