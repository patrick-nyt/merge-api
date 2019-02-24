package main

import (
	"context"
	"sort"
	"testing"
	"time"
)

func TestGatherReturnsMaxPosts(t *testing.T) {

	ctx := context.Background()
	service := &service{Config{}, twoPosts{}, threePosts{}}
	posts := service.gatherPosts(ctx, 5, true)
	if len(posts.Posts) != 5 {
		t.Errorf(
			"got %d posts expected %d",
			len(posts.Posts),
			5,
		)
	}
	if !sort.IsSorted(posts) {
		t.Error("got unsorted posts; expected sorted ascending")
	}
	if posts.Posts[0].Title != "Last Month" {
		t.Errorf("got %s first; expected %s",
			posts.Posts[0].Title,
			"Last Month",
		)
	}

	posts = service.gatherPosts(ctx, 3, false)
	if len(posts.Posts) != 3 {
		t.Errorf(
			"got %d posts expected %d",
			len(posts.Posts),
			3,
		)
	}
	if sort.IsSorted(posts) {
		t.Error("got sorted ascending posts; expected sorted descending")
	}
	if posts.Posts[0].Title != "Tomorrow" {
		t.Errorf("got %s first; expected %s",
			posts.Posts[0].Title,
			"Tomorrow",
		)
	}

}

type twoPosts struct{}

func (s twoPosts) Fetch(ctx context.Context, max int) []Post {
	today := time.Now()
	lastWeek := today.AddDate(0, 0, -7)
	return []Post{
		Post{
			Title:     "Today",
			CreatedAt: JSONTime{today},
		},
		Post{
			Title:     "Last Week",
			CreatedAt: JSONTime{lastWeek},
		},
	}
}

type threePosts struct{}

func (s threePosts) Fetch(ctx context.Context, max int) []Post {
	today := time.Now()
	yesterday := today.AddDate(0, 0, -1)
	tomorrow := today.AddDate(0, 0, 1)
	lastMonth := today.AddDate(0, -1, 0)
	return []Post{
		Post{
			Title:     "Yesterday",
			CreatedAt: JSONTime{yesterday},
		},
		Post{
			Title:     "Tomorrow",
			CreatedAt: JSONTime{tomorrow},
		},
		Post{
			Title:     "Last Month",
			CreatedAt: JSONTime{lastMonth},
		},
	}
}
