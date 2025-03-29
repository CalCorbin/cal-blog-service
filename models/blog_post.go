package models

import "gorm.io/gorm"

type BlogPost struct {
	gorm.Model
	Title   string
	Content string
	Author  string
}

func (b *BlogPost) Validate() bool {
	return b.Title != "" && b.Content != "" && b.Author != ""
}
