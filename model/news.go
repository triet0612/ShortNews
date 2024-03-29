package model

import (
	"encoding/xml"
)

type RSSNews struct {
	XMLName xml.Name   `xml:"rss"`
	Channel RSSChannel `xml:"channel"`
}

type RSSChannel struct {
	XMLName     xml.Name  `xml:"channel"`
	Title       string    `xml:"title"`
	Description string    `xml:"description"`
	PubDate     string    `xml:"pubDate"`
	Generator   string    `xml:"generator"`
	Link        string    `xml:"link"`
	Items       []RSSItem `xml:"item"`
}

type RSSImage struct {
	XMLName xml.Name `xml:"enclosure"`
	URL     string   `xml:"url,attr"`
}

type RSSItem struct {
	ArticleID string
	Link      string   `xml:"link"`
	Title     string   `xml:"title"`
	PubDate   string   `xml:"pubDate"`
	Image     RSSImage `xml:"enclosure"`
	Publisher string
	Summary   string
}

type NewsSource struct {
	PublisherID string
	Publisher   string
	Link        string
	Language    string
}
