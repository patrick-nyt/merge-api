package main

import (
	"context"
	"log"

	"cloud.google.com/go/datastore"
	"github.com/machinebox/graphql"
)

// Fetcher is an interface for something that returns
// a slice of Post from some source, given a max size.
type Fetcher interface {
	Fetch(context.Context, int) []Post
}

// remoteClient is a graphql client
type remoteClient struct {
	*graphql.Client
}

// NewRemoteClient makes a new graphql client
func NewRemoteClient(location string) *remoteClient {
	remote := &remoteClient{graphql.NewClient(location)}
	// uncomment below to log from graphql client
	// remote.Client.Log = func(s string) { log.Println(s) }
	return remote
}

// localClient is a client for a local db
type localClient struct {
	*datastore.Client
}

// NewRemoteClient makes a new client for the local db
func NewLocalClient(location string) *localClient {
	ctx := context.Background()
	local, err := datastore.NewClient(ctx, DB_PROJECT)
	if err != nil {
		log.Println("unable to start datastore client")
		return &localClient{}
	}
	return &localClient{local}
}

// allPostsQuery is the query used to retrieve a group of posts from
// the graphql source.
const allPostsQuery = `
	{
  		allPosts(count: $max) {
    		id
    		title
    		createdAt
    		author {
      			id
      			firstName
      			lastName
    		}
  		}
	}
`

// AllPosts is a group of posts. Used when unmarshaling
// the json reply from the graphql remote source.
type AllPosts struct {
	Data struct {
		Posts []Post `json:"allPosts"`
	} `json:"data"`
}

func (c *remoteClient) Fetch(ctx context.Context, max int) []Post {

	// prepare the graphql query
	req := graphql.NewRequest(allPostsQuery)
	req.Var("max", max)

	// execute the query
	var response AllPosts
	err := c.Run(ctx, req, &response)

	if err != nil {
		// log error and return no posts
		log.Println(err)
		return mockPosts(max) // todo: should return []Post{}; returning mock data since fakerql is not working
	}

	return response.Data.Posts
}

const (
	// todo: make these datastore params configurable
	DB_KIND      = "post"
	DB_NAMESPACE = "merge-test"
	DB_PROJECT   = "nyt-uhub-dev"
)

func (c *localClient) Fetch(ctx context.Context, max int) []Post {
	var data []Post
	q := datastore.NewQuery(DB_KIND).Namespace(DB_NAMESPACE).Limit(max)
	_, err := c.Client.GetAll(ctx, q, &data)

	if err != nil {
		// log error and return no posts
		log.Println("datastore query failed")
		return []Post{}
	}

	return data
}
