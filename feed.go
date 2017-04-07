package BlogggerBot

import (
	"time"

	"github.com/mmcdole/gofeed"
)

type FeedService interface {
	GetNewPosts(url string, since time.Time) ([]gofeed.Item, error)
	GetBlog(url string) (gofeed.Feed, error)
}
