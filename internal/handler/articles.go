package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"newscrapper/internal/model"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func (h *Handler) GetArticles(w http.ResponseWriter, r *http.Request) {
	limit := 1
	start := 0
	queryParams := r.URL.Query()
	var err error
	if limit, err = strconv.Atoi(queryParams.Get("limit")); err != nil {
		limit = 1
	}
	if start, err = strconv.Atoi(queryParams.Get("start")); err != nil {
		start = 0
	}
	articles := []model.Article{}
	rows, err := h.db.QueryContext(context.Background(),
		"SELECT * FROM Article ORDER BY datetime(PubDate) DESC LIMIT ? OFFSET ?",
		limit, start,
	)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "error get article", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		a := model.Article{}
		if err := rows.Scan(
			&a.ArticleID, &a.Link, &a.Title, &a.PubDate, &a.PublisherID, &a.Summary,
		); err != nil {
			slog.Error(err.Error())
			continue
		}
		articles = append(articles, a)
	}
	if len(articles) == 0 {
		http.Error(w, "no articles found", http.StatusNotFound)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(articles); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) GetArticleFilterID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := uuid.Validate(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	a := &model.Article{}
	row := h.db.QueryRowContext(context.Background(),
		"SELECT * FROM Article WHERE ArticleID=?", id)
	if err := row.Scan(
		&a.ArticleID, &a.Link, &a.Title, &a.PubDate, &a.PublisherID, &a.Summary,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "error no article", http.StatusNotFound)
			return
		}
		slog.Error(err.Error())
		http.Error(w, "error get articles", http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(a); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) GetArticleFullFilterID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := uuid.Validate(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	link := ""
	row := h.db.QueryRowContext(context.Background(),
		"SELECT Link FROM Article WHERE ArticleID=?", id)
	if err := row.Scan(&link); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "error no article", http.StatusNotFound)
			return
		}
		slog.Error(err.Error())
		http.Error(w, "error get articles", http.StatusInternalServerError)
		return
	}
	client := http.Client{Timeout: 10 * time.Second}
	res, err := client.Get(link)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "get article error", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()
	ans, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "error decode body", http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "text/html")
	w.Write(ans)
}

func (h *Handler) GetArticleThumbnail(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := uuid.Validate(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ans := []byte{}
	row := h.db.QueryRowContext(context.Background(),
		"SELECT Image FROM Thumbnail WHERE ArticleID=?", id)
	if err := row.Scan(&ans); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "no image", http.StatusNotFound)
			return
		}
		slog.Error(err.Error())
		http.Error(w, "error get thumbnail", http.StatusInternalServerError)
		return
	}
	if len(ans) == 0 {
		http.Error(w, "no image", http.StatusNotFound)
		return
	}
	w.Write(ans)
}

func (h *Handler) GetArticleAudio(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := uuid.Validate(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ans := []byte{}
	row := h.db.QueryRowContext(context.Background(),
		`SELECT Audio FROM ArticleAudio WHERE ArticleID=?`, id,
	)
	if err := row.Scan(&ans); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "no image", http.StatusNotFound)
			return
		}
		slog.Error(err.Error())
		http.Error(w, "error getting article audio", http.StatusInternalServerError)
	}
	if len(ans) == 0 {
		http.Error(w, "no audio", http.StatusNotFound)
		return
	}
	w.Write(ans)
}

func (h *Handler) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := uuid.Validate(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := context.Background()
	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "error delete article", http.StatusInternalServerError)
		slog.Error(tx.Rollback().Error())
		return
	}
	if _, err := tx.ExecContext(ctx, "DELETE FROM ArticleAudio WHERE ArticleID=?", id); err != nil {
		slog.Error(err.Error())
		http.Error(w, "error delete audio", http.StatusInternalServerError)
		slog.Error(tx.Rollback().Error())
		return
	}
	if _, err := tx.ExecContext(ctx, "DELETE FROM Thumbnail WHERE ArticleID=?", id); err != nil {
		slog.Error(err.Error())
		http.Error(w, "error delete thumbnail", http.StatusInternalServerError)
		slog.Error(tx.Rollback().Error())
		return
	}
	if _, err := tx.ExecContext(ctx, "DELETE FROM Article WHERE ArticleID=?", id); err != nil {
		slog.Error(err.Error())
		http.Error(w, "error delete article", http.StatusInternalServerError)
		slog.Error(tx.Rollback().Error())
		return
	}
	if err := tx.Commit(); err != nil {
		slog.Error(err.Error())
		http.Error(w, "error delete article", http.StatusInternalServerError)
	}
}
