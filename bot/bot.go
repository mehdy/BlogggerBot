package bot

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/mehdy/BlogggerBot"
	"github.com/spf13/viper"
)

type Bot struct {
	b        *tgbotapi.BotAPI
	handlers map[string]func(tgbotapi.Update)
	BS       BlogggerBot.BlogService
	FS       BlogggerBot.FeedService
}

func NewBot() *Bot {
	bot, err := tgbotapi.NewBotAPI(viper.GetString("BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)

	b := Bot{b: bot, handlers: map[string]func(tgbotapi.Update){}}

	b.handlers["/new_blog"] = b.AddNewBlog
	b.handlers["/update_posts"] = b.UpdatePosts
	b.handlers["/send_updates"] = b.SendNewPosts

	return &b
}

func (b *Bot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates, err := b.b.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		if update.Message == nil {
			log.Print(update)
			continue
		}
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.From.UserName != viper.GetString("BOT_ADMIN") {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You are not authorized to run command")
			b.b.Send(msg)
			continue
		}

		parsed := strings.SplitN(update.Message.Text, " ", 2)
		if handler, ok := b.handlers[parsed[0]]; ok {
			go handler(update)
		} else {
			go b.defaultHandler(update)
		}

	}
}

func (b *Bot) defaultHandler(u tgbotapi.Update) {
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Command not found")
	b.b.Send(msg)
}

func (b *Bot) SendNewPosts(u tgbotapi.Update) {
	posts, err := b.BS.GetNewPosts()
	if err != nil {
		log.Panicln(err)
	}
	for _, post := range posts {
		text := fmt.Sprintf(viper.GetString("BOT_MESSAGE_TEMPLATE"), post.Blog.Title, post.Title, post.URL)
		msg := tgbotapi.NewMessageToChannel(viper.GetString("BOT_CHANNEL"), text)
		b.b.Send(msg)
		if err := b.BS.Notify(&post); err != nil {
			log.Println(err)
		}
	}
}

func (b *Bot) UpdatePosts(u tgbotapi.Update) {
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, "starting to update the posts...")
	b.b.Send(msg)

	blogs, err := b.BS.GetBlogs()
	if err != nil {
		log.Print(err)
		msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Unable to get blogs :(")
		b.b.Send(msg)
		return
	}
	for _, blog := range blogs {
		feeds, err := b.FS.GetNewPosts(blog.FeedURL, blog.LastUpdatedAt)
		if err != nil {
			log.Println(err)
		}

		for _, feed := range feeds {
			p := BlogggerBot.Post{
				BlogID:      blog.ID,
				Title:       feed.Title,
				PublishedAt: *feed.PublishedParsed,
				URL:         feed.Link,
				GUID:        feed.GUID,
			}
			b.BS.CreatePost(&p)
		}

		b.BS.UpdatedBlog(&blog)
	}

	msg = tgbotapi.NewMessage(u.Message.Chat.ID, "Updated Successfully")
	b.b.Send(msg)
}

func (b *Bot) AddNewBlog(u tgbotapi.Update) {
	var msg tgbotapi.MessageConfig
	parsed := strings.SplitN(u.Message.Text, " ", 2)
	if len(parsed) != 2 {
		msg = tgbotapi.NewMessage(u.Message.Chat.ID, "No URL to add. command should be like /add_blog URL")
	}
	feed, err := b.FS.GetBlog(parsed[1])
	if err != nil {
		log.Println(err)
		msg = tgbotapi.NewMessage(u.Message.Chat.ID, "Unable to add this feed")
	} else {
		blog := BlogggerBot.Blog{
			Title:         feed.Title,
			Description:   feed.Description,
			Language:      feed.Language,
			FeedURL:       feed.Link,
			LastUpdatedAt: time.Now(),
		}
		if err := b.BS.CreateBlog(&blog); err != nil {
			log.Println(err)
			msg = tgbotapi.NewMessage(u.Message.Chat.ID, "Unable to add this feed")
		} else {
			msg = tgbotapi.NewMessage(u.Message.Chat.ID, "Blog has been added successfully")
		}
	}
	b.b.Send(msg)
}
