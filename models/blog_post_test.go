package models

import (
	"testing"
)

func TestBlogPostCreation(t *testing.T) {
	// Create a new blog post
	post := BlogPost{
		Title:    "Cowboy Bebop",
		Content:  "See you later space cowboy",
		AuthorID: "Spike Spiegel",
	}

	// Verify fields are set correctly
	if post.Title != "Cowboy Bebop" {
		t.Errorf("Title should be 'Cowboy Bebop', got '%s'", post.Title)
	}

	if post.Content != "See you later space cowboy" {
		t.Errorf("Expected Content to be 'See you later space cowboy', got '%s'", post.Content)
	}

	if post.AuthorID != "Spike Spiegel" {
		t.Errorf("Expected AuthorID to be 'Spike Spiegel', got '%s'", post.AuthorID)
	}
}

func TestBlogPostValid(t *testing.T) {
	post := BlogPost{
		Title:    "Cowboy Bebop",
		Content:  "See you later space cowboy",
		AuthorID: "Spike Spiegel",
	}

	if !post.Validate() {
		t.Errorf("BlogPost should be valid")
	}
}

func TestBlogPostInvalid(t *testing.T) {
	post := BlogPost{
		Title:    "",
		Content:  "See you later space cowboy",
		AuthorID: "Spike Spiegel",
	}

	if post.Validate() {
		t.Errorf("BlogPost should be invalid")
	}
}
