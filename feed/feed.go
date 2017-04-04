package main

import "github.com/mmcdole/gofeed"

type FeedReader struct {
	f *gofeed.Parser
}

func NewFeedReader() *FeedReader {
	var fr FeedReader
	fr.f = gofeed.NewParser()
	return &fr
}

func (f *FeedReader) GetNewPosts(url string) error {
	feed, err := f.f.ParseURL(url)
	if err != nil {
		return nil
	}
}
