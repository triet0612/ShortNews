package handler

import (
	"net/http"
	"newscrapper/internal/config"
)

type Handler struct {
	DI config.DI
}

func (h *Handler) Mount(mux *http.ServeMux) {
	mux.HandleFunc("GET /", h.Home)

	mux.HandleFunc("GET /api/rss", h.GetRssSource)
	mux.HandleFunc("POST /api/rss", h.CreateRssSource)
	mux.HandleFunc("PUT /api/rss", h.UpdateRssSource)
	mux.HandleFunc("DELETE /api/rss", h.DeleteRssSource)

	mux.HandleFunc("GET /api/articles", h.GetArticles)
	mux.HandleFunc("GET /api/articles/{uuid}", h.GetArticlesFilterID)
	mux.HandleFunc("GET /api/articles/thumbnail/{uuid}", h.GetArticleThumbnail)
	mux.HandleFunc("PUT /api/articles/{uuid}", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("DELETE /api/articles/{uuid}", func(w http.ResponseWriter, r *http.Request) {})
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("hello world"))
}
