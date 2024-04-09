package handler

import (
	"net/http"
	"newscrapper/internal/config"
	"newscrapper/internal/db"
	"os"
)

type Handler struct {
	db     *db.DBService
	clock  *config.Clock
	signal chan struct{}
}

type ExtendWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *ExtendWriter) Write(b []byte) (int, error) {
	if w.statusCode == http.StatusNotFound {
		return len(b), nil
	}
	if w.statusCode != 0 {
		w.WriteHeader(w.statusCode)
	}
	return w.ResponseWriter.Write(b)

}

func (w *ExtendWriter) WriteHeader(statusCode int) {
	if statusCode >= 300 && statusCode < 400 {
		w.ResponseWriter.WriteHeader(statusCode)
		return
	}
	w.statusCode = statusCode
}

func NewHandler(db *db.DBService, clock *config.Clock, signal chan struct{}) *Handler {
	return &Handler{db: db, clock: clock, signal: signal}
}

func (h *Handler) Mount(mux *http.ServeMux) {
	fsDir := os.DirFS("./build")
	fileServer := http.FileServer(http.FS(fsDir))

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		ex := &ExtendWriter{ResponseWriter: w}
		fileServer.ServeHTTP(ex, r)
		if ex.statusCode == http.StatusNotFound {
			r.URL.Path = "/"
			w.Header().Set("Content-Type", "text/html")
			fileServer.ServeHTTP(w, r)
		}
	})

	mux.HandleFunc("GET /api/rss", h.GetRssSource)
	mux.HandleFunc("POST /api/rss", h.CreateRssSource)
	mux.HandleFunc("DELETE /api/rss/{id}", h.DeleteRssSource)

	mux.HandleFunc("GET /api/articles", h.GetArticle)
	mux.HandleFunc("GET /api/articles/random", h.GetRandomArticle)
	mux.HandleFunc("GET /api/articles/{id}", h.GetArticleFilterID)
	mux.HandleFunc("DELETE /api/articles/{id}", h.DeleteArticle)

	mux.HandleFunc("GET /api/articles/thumbnail/{id}", h.GetArticleThumbnail)
	mux.HandleFunc("GET /api/articles/audio/{id}", h.GetArticleAudio)
}
