package models

import "gorm.io/gorm"

type BlogPost struct {
	gorm.Model
	Title    string `json:"title"`
	Content  string `json:"content"`
	AuthorID uint   `json:"author_id"`
}

func (b *BlogPost) Validate() bool {
	return b.Title != "" && b.Content != "" && b.AuthorID != 0
}
