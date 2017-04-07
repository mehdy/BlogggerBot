package main

import (
	"log"

	"github.com/mehdy/BlogggerBot/bot"
	"github.com/mehdy/BlogggerBot/feed"
	"github.com/mehdy/BlogggerBot/models"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("bloggger")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	app := bot.NewBot()

	db := models.OpenDB()
	defer db.Close()

	app.BS = &models.BlogService{db}
	app.FS = feed.NewFeedReader()

	app.Run()
}
