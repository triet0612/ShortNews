package handler

import (
	"net/http"
	"newscrapper/internal/config"
	"newscrapper/internal/db"
)

type Handler struct {
	db    *db.DBService
	clock *config.Clock
}

func NewHandler(db *db.DBService, clock *config.Clock) *Handler {
	return &Handler{db: db, clock: clock}
}

func (h *Handler) Mount(mux *http.ServeMux) {
	mux.HandleFunc("GET /", h.Home)

	mux.HandleFunc("GET /api/rss", h.GetRssSource)
	mux.HandleFunc("POST /api/rss", h.CreateRssSource)
	mux.HandleFunc("PUT /api/rss", h.UpdateRssSource)
	mux.HandleFunc("DELETE /api/rss/{id}", h.DeleteRssSource)

	mux.HandleFunc("GET /api/articles", h.GetArticles)
	mux.HandleFunc("GET /api/articles/{id}", h.GetArticleFilterID)
	mux.HandleFunc("GET /api/articles/full/{id}", h.GetArticleFullFilterID)
	mux.HandleFunc("GET /api/articles/thumbnail/{id}", h.GetArticleThumbnail)
	mux.HandleFunc("GET /api/articles/audio/{id}", h.GetArticleAudio)
	mux.HandleFunc("DELETE /api/articles/{id}", h.DeleteArticle)
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("hello world"))
}
