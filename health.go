package main

import (
	"log"
	"net/http"
)

// health always returns 200, {"Status":"OK"} to the caller.
func (s *service) health(w http.ResponseWriter, r *http.Request) {

	log.Println("Request: health")

	// as part of this health check, may want to check service dependencies

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"Status":"OK"}`))

}
