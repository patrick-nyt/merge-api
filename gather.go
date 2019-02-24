package main

import (
	"context"
	"sort"
	"sync"
)

// Post is a single post returned from a source
type Post struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Author    Author   `json:"author"`
	CreatedAt JSONTime `json:"createdAt"`
}

// Author is the author of a post
type Author struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// Posts is a group of posts. Used when marshaling the json response.
type Posts struct {
	Posts []Post `json:"posts"`
}

// gatherPosts retrieves and merges at most max posts from the configured sources,
// sorted in chronlogical order if ascending is true, otherwise reverse
// chronological order.
func (s *service) gatherPosts(ctx context.Context, max int, ascending bool) Posts {

	var remotePosts, localPosts, gatheredPosts []Post

	var wg sync.WaitGroup

	// get posts from the remote source asynchronously
	wg.Add(1)
	go func() {
		defer wg.Done()
		remotePosts = s.remote.Fetch(ctx, max)
	}()

	// get posts from the local source asynchronously
	wg.Add(1)
	go func() {
		defer wg.Done()
		localPosts = s.local.Fetch(ctx, max)
	}()

	// wait for the two async calls to finish
	wg.Wait()

	// Combine the two lists, sort the result
	// and return at most max items.
	//
	// Optimization are left until later, if they needed.
	// For instance, if we know that each list is already sorted,
	// then we can merge the lists instead of sorting.
	// Or memory could be allocated up front instead of
	// relying of go's memory management.
	gatheredPosts = localPosts
	gatheredPosts = append(gatheredPosts, remotePosts...)
	posts := Posts{gatheredPosts}
	if ascending {
		sort.Sort(posts)
	} else {
		sort.Sort(sort.Reverse(posts))
	}

	// return at most max items
	if len(posts.Posts) > max {
		posts.Posts = posts.Posts[:max]
	}

	// Note that if both async calls fail,
	// this will simply return no postings.
	// Might want to handle that case diferently
	// to be able to return an http error to
	// the client.
	return posts
}

// The methods Len(), Swap() and Less() are added
// so Posts will implement the sort inferface.
// The sort is ordered by the createdAt time.

func (p Posts) Len() int {
	return len(p.Posts)
}
func (p Posts) Swap(i, j int) {
	p.Posts[i], p.Posts[j] = p.Posts[j], p.Posts[i]
}
func (p Posts) Less(i, j int) bool {
	return p.Posts[i].CreatedAt.Time.Before(p.Posts[j].CreatedAt.Time)
}
