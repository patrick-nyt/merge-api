package main

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// Posts returns a slice of n mock posts with random data.
// Included because fakerql is currently not working, otherwise
// should just use for testing.
func mockPosts(n int) []Post {
	result := make([]Post, n)
	for r := range result {
		result[r] = Post{
			rand5("ID"),
			rand5("Title"),
			Author{rand5("ID"),
				rand5("FirstName"),
				rand5("LastName")},
			JSONTime{randTime()},
		}
	}
	return result
}

// randTime returns time a random few hours and days in the future.
func randTime() time.Time {
	return time.Now().AddDate(0, rand.Intn(16), rand.Intn(64))
}

// rand5 adds 5 random letters to a string
func rand5(s string) string { return s + "_" + randString(5) }

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// randString returns a random string with n letters.
func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
