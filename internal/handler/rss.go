package handler

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"newscrapper/internal/model"

	"github.com/google/uuid"
)

func (h *Handler) GetRssSource(w http.ResponseWriter, r *http.Request) {
	src, err := h.db.ReadSourceRSS()
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
	lang, ok := newRssSource["language"]
	if !ok {
		http.Error(w, "no language found", http.StatusInternalServerError)
		return
	}
	url, err := url.ParseRequestURI(rssLink)
	if err != nil {
		http.Error(w, "POST /api/rss error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	src := &model.NewsSource{
		PublisherID: uuid.NewString(),
		Link:        url.String(),
		Language:    lang,
	}
	if _, err := h.db.ExecContext(context.Background(),
		"INSERT INTO NewsSource VALUES (?, ?, ?)",
		src.PublisherID, src.Link, src.Language,
	); err != nil {
		slog.Error(err.Error())
		http.Error(w, "error create rss source", http.StatusInternalServerError)
		return
	}
	h.signal <- struct{}{}
}

func (h *Handler) DeleteRssSource(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")
	if err := uuid.Validate(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := context.Background()
	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "error prepare delete news source", http.StatusInternalServerError)
		return
	}
	if _, err := tx.ExecContext(ctx,
		`DELETE FROM Thumbnail WHERE ArticleID IN (
SELECT ArticleID FROM Article WHERE PublisherID = ?);`, id); err != nil {
		slog.Error(err.Error())
		tx.Rollback()
		http.Error(w, "error delete thumbnail", http.StatusInternalServerError)
		return
	}
	if _, err := tx.ExecContext(ctx,
		`DELETE FROM ArticleAudio WHERE ArticleID IN (
SELECT ArticleID FROM Article WHERE PublisherID=?);`, id); err != nil {
		slog.Error(err.Error())
		tx.Rollback()
		http.Error(w, "error delete audio", http.StatusInternalServerError)
		return
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM Article WHERE PublisherID=?;`, id); err != nil {
		slog.Error(err.Error())
		tx.Rollback()
		http.Error(w, "error delete article", http.StatusInternalServerError)
		return
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM NewsSource WHERE PublisherID=?;`, id); err != nil {
		slog.Error(err.Error())
		tx.Rollback()
		http.Error(w, "error delete rss", http.StatusInternalServerError)
		return
	}
	if err := tx.Commit(); err != nil {
		slog.Error(err.Error())
		tx.Rollback()
		http.Error(w, "error delete rss", http.StatusInternalServerError)
	}
	h.signal <- struct{}{}
}
