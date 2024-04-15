package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"newscrapper/internal/model"
	"strconv"

	"github.com/google/uuid"
)

func (h *Handler) GetArticle(w http.ResponseWriter, r *http.Request) {
	limit := 1
	start := 0
	summary, audio, pubid := "", "", ""
	queryParams := r.URL.Query()
	var err error
	if limit, err = strconv.Atoi(queryParams.Get("limit")); err != nil {
		limit = 1
	}
	if start, err = strconv.Atoi(queryParams.Get("start")); err != nil {
		start = 0
	}
	summary = queryParams.Get("summary")
	audio = queryParams.Get("audio")
	pubid = queryParams.Get("PublisherID")
	query := `SELECT * FROM Article WHERE 1=1`
	if audio == "true" {
		query += ` AND ArticleID IN (SELECT ArticleID FROM ArticleAudio WHERE Audio != "")`
	}
	if summary == "true" {
		query += ` AND Summary != ""`
	}
	if pubid != "" {
		query += ` AND PublisherID == @pubid`
	} else {
		query += ` AND PublisherID != @pubid`
	}
	query += ` ORDER BY datetime(PubDate) DESC LIMIT @limit OFFSET @offset`
	articles := []model.Article{}
	ctx := context.Background()
	stmt, err := h.db.PrepareContext(ctx, query)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "error getting articles", http.StatusInternalServerError)
		return
	}
	rows, err := stmt.QueryContext(ctx, sql.Named("pubid", pubid), sql.Named("limit", limit), sql.Named("offset", start))
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

func (h *Handler) GetArticleThumbnail(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := uuid.Validate(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ans := ""
	row := h.db.QueryRowContext(context.Background(),
		"SELECT URL FROM Thumbnail WHERE ArticleID=?", id)
	if err := row.Scan(&ans); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "no image", http.StatusNotFound)
			return
		}
		slog.Error(err.Error())
		http.Error(w, "error get thumbnail", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, ans, http.StatusTemporaryRedirect)
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
	w.Header().Add("audio", "wav")
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
