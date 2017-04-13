package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/mehdy/BlogggerBot"
	"github.com/spf13/viper"
)

type Bot struct {
	b        *tgbotapi.BotAPI
	hc       *http.Client
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

	b := Bot{
		b:        bot,
		hc:       &http.Client{Timeout: 3 * time.Second},
		handlers: map[string]func(tgbotapi.Update){},
	}

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
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, "starting to send updated posts...")
	b.b.Send(msg)
	posts, err := b.BS.GetNewPosts()
	if err != nil {
		log.Panicln(err)
		msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Unable to get updated posts")
		b.b.Send(msg)
	}
	for _, post := range posts {
		text := fmt.Sprintf(viper.GetString("BOT_MESSAGE_TEMPLATE"), post.Author, post.Title, post.ShortURL)
		msg := tgbotapi.NewMessageToChannel(viper.GetString("BOT_CHANNEL"), text)
		b.b.Send(msg)
		if err := b.BS.Notify(&post); err != nil {
			log.Println(err)
		}
	}
	msg = tgbotapi.NewMessage(u.Message.Chat.ID, "sent all updated posts")
	b.b.Send(msg)

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
			var shortURL string
			body := bytes.NewBuffer([]byte(fmt.Sprintf(`{"longUrl": "%s"}`, feed.Link)))
			reqUrl := fmt.Sprintf("https://www.googleapis.com/urlshortener/v1/url?key=%s",
				viper.GetString("API_TOKEN"))
			req, err := http.NewRequest("POST", reqUrl, body)
			if err == nil {
				req.Header.Set("Content-Type", "application/json")
				resp, err := b.hc.Do(req)
				if err == nil {
					var res map[string]string
					if err := json.NewDecoder(resp.Body).Decode(&res); err == nil {
						if id, ok := res["id"]; ok {
							shortURL = id
						}
					}
				}
			}

			if shortURL == "" {
				shortURL = feed.Link
			}
			p := BlogggerBot.Post{
				BlogID:      blog.ID,
				Title:       feed.Title,
				PublishedAt: *feed.PublishedParsed,
				URL:         feed.Link,
				ShortURL:    shortURL,
				GUID:        feed.GUID,
			}

			if feed.Author != nil && feed.Author.Name != "" {
				p.Author = feed.Author.Name
			} else {
				p.Author = blog.Title
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
			FeedURL:       parsed[1],
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
