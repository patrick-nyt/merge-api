package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// DEFAULT_MAX is the number of posts to return when this is not specified.
const DEFAULT_MAX = 25

// MAX_MAX is the largest number of posts that can be returned.
// This particular number is arbitrary-- probably the limit should
// be determined by the number of posts that can be queried from
// the graphql and database sources, or maybe the number of posts
// that the client of this api can consume. Or perhaps there should
// be no limit and the posts should be streamed to the client?
const MAX_MAX = 10000

// DEFAULT_ORDER is the default sort order for the returned posts.
const DEFAULT_ORDER = "desc"

// posts is an http handler which retrieves postings by merging two sources.
// Query params are "max" and "order". Returns at most max results (default: 25),
// ordered either ascending ("asc") or descending ("desc") -- (default: "asc").
// Returns 200 on success, 400 if validation fails, and 500 if something unexpted happens.
func (s *service) posts(w http.ResponseWriter, r *http.Request) {

	var err error

	// get query params and do some basic validation
	max := DEFAULT_MAX
	maxParam := r.URL.Query().Get("max")
	if maxParam != "" {
		max, err = strconv.Atoi(maxParam)
		if err != nil || max < 1 || max > MAX_MAX {
			http.Error(w, "invalid max param", http.StatusBadRequest)
			return
		}
	}
	order := r.URL.Query().Get("order")
	if order == "" {
		order = DEFAULT_ORDER
	}
	if order != "asc" && order != "desc" {
		http.Error(w, "invalid order param", http.StatusBadRequest)
		return
	}

	log.Printf("Request posts: max %d, order %s\n", max, order)

	// retrieve a group of at most max postings
	ctx := context.Background()
	postings := s.gatherPosts(ctx, max, order == "asc")

	// return result in json
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(postings)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
