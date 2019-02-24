package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

// a few tests for request validation and error handling on the posts endpoint.

func TestBadOrderInPostsRequest(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	q := req.URL.Query()
	q.Add("order", "badorder")
	req.URL.RawQuery = q.Encode()

	rr := httptest.NewRecorder()

	service := service{config: Config{}}
	handler := http.HandlerFunc(service.posts)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf(
			"unexpected status: got (%v) want (%v)",
			status,
			http.StatusBadRequest,
		)
	}
}

func TestNegativeMaxInPostsRequest(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	q := req.URL.Query()
	q.Add("max", "-2")
	req.URL.RawQuery = q.Encode()

	rr := httptest.NewRecorder()

	service := service{config: Config{}}
	handler := http.HandlerFunc(service.posts)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf(
			"unexpected status: got (%v) want (%v)",
			status,
			http.StatusBadRequest,
		)
	}
}

func TestValidPostsRequest(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	q := req.URL.Query()
	q.Add("order", "desc")
	q.Add("max", "45")
	req.URL.RawQuery = q.Encode()

	rr := httptest.NewRecorder()

	log.SetOutput(ioutil.Discard) // turn off logging

	service := &service{Config{}, emptySource{}, emptySource{}}
	handler := http.HandlerFunc(service.posts)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf(
			"unexpected status: got (%v) want (%v)",
			status,
			http.StatusOK,
		)
	}
}

type emptySource struct{}

func (s emptySource) Fetch(ctx context.Context, max int) []Post {
	return []Post{}
}
