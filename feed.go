package main

import "github.com/mmcdole/gofeed"

type FeedService interface {
	GetNewPosts(url string) ([]gofeed.Item, error)
	GetBlog(url string) (gofeed.Feed, error)
}
