package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Env          string `envconfig:"ENV" default:"dev"`
	Port         string `envconfig:"PORT" default:"8080"`
	RemoteSource string `envconfig:"MERGEAPI_REMOTE" default:"https://fakerql.com/graphql"`
	LocalSource  string `envconfig:"MERGEAPI_LOCAL" default:"https://localhost:9999"`
}

func main() {

	log.Printf("Initializing merge api")

	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Panicln(err)
	}
	log.Printf("Config: %+v", config)

	service := NewService(config)
	router := mux.NewRouter()
	router.HandleFunc("/health", service.health).Methods("GET")
	router.HandleFunc("/posts", service.posts).Methods("GET")

	log.Printf("Listening on port %s", config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", config.Port), router))
}
