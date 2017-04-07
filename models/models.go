package models

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/mehdy/BlogggerBot"
)

func OpenDB() *gorm.DB {
	var err error
	var db *gorm.DB

	if db, err = gorm.Open("sqlite3", "bloggger.db"); err != nil {
		log.Fatal(err)
	}

	if err = db.DB().Ping(); err != nil {
		log.Fatal(err)
	}

	// Migrations
	db.Debug().AutoMigrate(&BlogggerBot.Blog{}, &BlogggerBot.Post{})
	return db
}
