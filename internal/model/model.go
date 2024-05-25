package model

import (
	"time"
)

type NewsSource struct {
	PublisherID string
	Link        string
	Language    string
	VoiceType   string
	Ext         map[string]interface{}
}

type Article struct {
	ArticleID   string
	Link        string
	Title       string
	PubDate     *time.Time
	PublisherID string
	Summary     string
	Ext         map[string]interface{}
}

type ArticleThumbnail struct {
	ArticleID string
	URL       string
	Image     *[]byte
	Ext       map[string]interface{}
}
