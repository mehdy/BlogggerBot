package main

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Blog struct {
	gorm.Model
	Title         string    `gorm:"SIZE:512;NOT NULL"`
	Description   string    `gorm:"SIZE:2048"`
	URL           string    `gorm:"SIZE:1024;NOT NULL"`
	FeedURL       string    `gorm:"SIZE:1024;NOT NULL"`
	Language      string    `gorm:"SIZE:5"`
	LastUpdatedAt time.Time `gorm:"NOT NULL"`
	Posts         []Post    `gorm:"ForeignKey:BlogID"`
}

type Post struct {
	gorm.Model
	BlogID      int       `gorm:"NOT NULL"`
	Author      string    `gorm:"SIZE:256;NOT NULL"`
	Title       string    `gorm:"SIZE:2048;NOT NULL"`
	Summary     string    `gorm:"NOT NULL"`
	PublishedAt time.Time `gorm:"NOT NULL"`
	URL         string    `gorm:"NOT NULL"`
	GUID        string    `gorm:"NOT NULL;UNIQUE"`
	Notified    bool      `gorm:"NOT NULL;DEFAULT:false"`
}

type BlogService interface {
	GetBlog(int) (Blog, error)
	GetBlogs() ([]Blog, error)
	CreateBlog(*Blog) error
	CreatePost(*Post) error
	GetNewPosts() ([]Post, error)
	Notify(*Post) error
}
