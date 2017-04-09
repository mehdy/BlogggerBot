package models

import (
	"time"

	"github.com/jinzhu/gorm"

	"github.com/mehdy/BlogggerBot"
)

type BlogService struct {
	DB *gorm.DB
}

func (s *BlogService) GetBlog(id uint) (BlogggerBot.Blog, error) {
	var b BlogggerBot.Blog
	err := s.DB.First(&b, id).Error
	return b, err
}

func (s *BlogService) GetBlogs() ([]BlogggerBot.Blog, error) {
	var blogs []BlogggerBot.Blog
	err := s.DB.Find(&blogs).Error
	return blogs, err
}

func (s *BlogService) CreateBlog(b *BlogggerBot.Blog) error {
	return s.DB.Create(b).Error
}

func (s *BlogService) UpdatedBlog(b *BlogggerBot.Blog) error {
	return s.DB.Model(b).Update("last_updated_at", time.Now()).Error
}

func (s *BlogService) CreatePost(p *BlogggerBot.Post) error {
	return s.DB.Create(p).Error
}

func (s *BlogService) GetNewPosts() ([]BlogggerBot.Post, error) {
	var posts []BlogggerBot.Post
	err := s.DB.Not("notified", true).Find(&posts).Error
	return posts, err
}

func (s *BlogService) Notify(p *BlogggerBot.Post) error {
	return s.DB.Model(p).Update("notified", true).Error
}
