package bot

import (
	"fmt"
	"log"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/mehdy/BlogggerBot"
	"github.com/spf13/viper"
)

type Bot struct {
	b  *tgbotapi.BotAPI
	BS BlogggerBot.BlogSerivce
	FS BlogggerBot.FeedService
}

func NewBot() {
	bot, err := tgbotapi.NewBotAPI(viper.GetString("BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)
}

func (b *Bot) SendNewPosts() {
	posts, err := b.BS.GetNewPosts()
	if err != nil {
		log.Panicln(err)
	}
	for _, post := range posts {
		text := fmt.Sprintf(viper.GetString("BOT_MESSAGE_TEMPLATE"), post.Author, post.Title, post.URL)
		msg := tgbotapi.NewMessageToChannel(viper.GetString("BOT_CHANNEL"), text)
		b.b.Send(msg)
		if err := b.BS.Notify(&post); err != nil {
			log.Println(err)
		}
	}
}

func (b *Bot) UpdatePosts() {
	// TODO: reply in bot chat
	msg := tgbotapi.NewMessage("starting to update the posts...")
	b.b.Send(msg)

	blogs, err := b.BS.GetBlogs()
	for _, blog := range blogs {
		feeds, err := b.FS.GetNewPosts(blog.FeedURL, blog.LastUpdatedAt)
		if err != nil {
			log.Println(err)
		}

		for _, feed := range feeds {
			p := BlogggerBot.Post{
				BlogID:      blog.ID,
				Author:      feed.Author,
				Title:       feed.Title,
				Summary:     feed.Summary,
				PublishedAt: *feed.PublishedParsed,
				URL:         feed.Link,
				GUID:        feed.GUID,
			}
			b.BS.CreatePost(p)
		}
	}

	// TODO: reply in bot chat
	msg := tgbotapi.NewMessage("Updated Successfully")
	b.b.Send(msg)
}

func (b *Bot) AddNewBlog() {
	// TODO: receive a URL
	url := ""
	var msg tgbotapi.MessageConfig
	feed, err := b.FS.GetBlog(url)
	if err != nil {
		log.Println(err)
		// TODO: reply in bot chat
		msg = tgbotapi.NewMessage("Unable to add this feed")
	} else {
		blog := BlogggerBot.Blog{
			Title:         feed.Title,
			Description:   feed.Description,
			Language:      feed.Language,
			URL:           feed.Link,
			FeedURL:       feed.FeedLink,
			LastUpdatedAt: time.Now(),
		}
		if err := b.BS.CreateBlog(&blog); err != nil {
			log.Println(err)
			// TODO: reply in bot chat
			msg = tgbotapi.NewMessage("Unable to add this feed")
		} else {
			// TODO: reply in bot chat
			msg = tgbotapi.NewMessage("Blog has been added successfully")
		}
	}
	b.b.Send(msg)
}
