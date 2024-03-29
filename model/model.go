package model

import (
	"time"
)

type NewsSource struct {
	PublisherID string
	Publisher   string
	Link        string
	Language    string
}

type Article struct {
	ArticleID   string
	Link        string
	Title       string
	PubDate     *time.Time
	PublisherID string
	Summary     string
}

type ArticleThumbnail struct {
	ArticleID string
	URL       string
	Image     *[]byte
}
