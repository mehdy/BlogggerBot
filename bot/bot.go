package bot

import (
	"fmt"
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/viper"
)

type Bot struct {
	b *tgbotapi.BotAPI
}

func NewBot() {
	bot, err := tgbotapi.NewBotAPI(viper.GetString("BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)
}

func (b *Bot) NewBlogPost() {
	text := fmt.Sprintf(`📢 پست جدید «%s» در مورد «%s» بخونید و براش کامنت بذارید.

✅ %s`, author, title, link)
	msg := tgbotapi.NewMessageToChannel("", text)

	b.b.Send(msg)
}
