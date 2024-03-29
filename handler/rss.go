package handler

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"newscrapper/model"
	"time"

	"github.com/google/uuid"
	"github.com/mmcdole/gofeed"
)

func (h *Handler) GetRssSource(w http.ResponseWriter, r *http.Request) {
	src, err := h.DI.DbCon.ReadSourceRSS()
	if err != nil {
		log.Println(err)
		http.Error(w, "GET /rss error: no source found", http.StatusNotFound)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(src); err != nil {
		log.Println(err)
		http.Error(w, "GET /rss error json parsing", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) CreateRssSource(w http.ResponseWriter, r *http.Request) {
	newRssSource := map[string]string{}
	if err := json.NewDecoder(r.Body).Decode(&newRssSource); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rssLink, ok := newRssSource["link"]
	if !ok {
		http.Error(w, "no link found", http.StatusInternalServerError)
		return
	}
	url, err := url.ParseRequestURI(rssLink)
	if err != nil {
		http.Error(w, "POST /api/rss error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	fp := gofeed.NewParser()
	feed, err := fp.ParseURLWithContext(url.String(), ctx)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "cannot fetch feed", http.StatusInternalServerError)
	}

	rssSource := &model.NewsSource{
		PublisherID: uuid.NewString(),
		Publisher:   feed.Title,
		Link:        url.String(),
		Language:    feed.Language,
	}
	if err = h.DI.DbCon.InsertRSSsource(*rssSource); err != nil {
		slog.Error(err.Error())
		http.Error(w, "error create rss source", http.StatusInternalServerError)
		return
	}
	h.DI.Clock.Sync()
}

func (h *Handler) UpdateRssSource(w http.ResponseWriter, r *http.Request) {
	// TODO: Update RSS source by ID and Sync RSS after
}

func (h *Handler) DeleteRssSource(w http.ResponseWriter, r *http.Request) {
	// TODO: Delete RSS source by ID and related articles
}
