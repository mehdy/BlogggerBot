package feed

import (
	"time"

	"github.com/mmcdole/gofeed"
)

type FeedReader struct {
	f *gofeed.Parser
}

func NewFeedReader() *FeedReader {
	var fr FeedReader
	fr.f = gofeed.NewParser()
	return &fr
}

func (f *FeedReader) GetBlog(url string) (gofeed.Feed, error) {
	feed, err := f.f.ParseURL(url)
	return *feed, err
}

func (f *FeedReader) GetNewPosts(url string, since time.Time) ([]gofeed.Item, error) {
	feed, err := f.f.ParseURL(url)
	if err != nil {
		return []gofeed.Item{}, err
	}
	var posts []gofeed.Item
	for _, item := range feed.Items {
		if item.PublishedParsed != nil && item.PublishedParsed.After(since) {
			posts = append(posts, *item)
		}
	}
	return posts, nil
}
