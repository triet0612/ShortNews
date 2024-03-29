package handler

import (
	"encoding/json"
	"encoding/xml"
	"log"
	"net/http"
	"net/url"
	"newscrapper/internal/config"
	"newscrapper/model"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type Handler struct {
	DI config.DI
}

func (h *Handler) Mount(mux *http.ServeMux) {
	mux.HandleFunc("GET /", h.Home)

	mux.HandleFunc("GET /api/rss", h.GetRssSource)
	mux.HandleFunc("POST /api/rss", h.CreateRssSource)
	mux.HandleFunc("PUT /api/rss", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("DELETE /api/rss", func(w http.ResponseWriter, r *http.Request) {})

	mux.HandleFunc("GET /api/articles", h.GetArticles)
	mux.HandleFunc("GET /api/articles/{uuid}", h.GetArticlesFilterID)
	mux.HandleFunc("GET /api/articles/thumbnail/{uuid}", h.GetArticleThumbnail)
	mux.HandleFunc("PUT /api/aritcles/{uuid}", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("DELETE /api/articles/{uuid}", func(w http.ResponseWriter, r *http.Request) {})
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("hello world"))
}

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
		http.Error(w, "POST /api/rss error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	rssLink, ok := newRssSource["link"]
	if !ok {
		http.Error(w, "POST /api/rss error: no link found", http.StatusInternalServerError)
		return
	}
	url, err := url.ParseRequestURI(rssLink)
	if err != nil {
		http.Error(w, "POST /api/rss error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	client := http.Client{Timeout: 10 * time.Second}
	res, err := client.Get(url.String())
	if err != nil {
		http.Error(w, "POST /api/rss error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	rss := model.RSSNews{}
	if err := xml.NewDecoder(res.Body).Decode(&rss); err != nil {
		http.Error(w, "POST /api/rss error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	rssSource := &model.NewsSource{
		PublisherID: uuid.NewString(),
		Publisher:   rss.Channel.Title,
		Link:        rssLink,
		Language:    "",
	}
	if err = h.DI.DbCon.InsertRSSsource(*rssSource); err != nil {
		log.Println("CreateNewRssSource err: ", err)
		http.Error(w, "CreateNewRssSource error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	h.DI.Clock.Sync()
}

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
	articles, err := h.DI.DbCon.ReadArticle(uint(limit), uint(start))
	if err != nil {
		http.Error(w, "GET /api/articles: ReadArticle", http.StatusInternalServerError)
		return
	}
	if len(*articles) == 0 {
		http.Error(w, "no articles found", http.StatusNotFound)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(articles); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) GetArticlesFilterID(w http.ResponseWriter, r *http.Request) {
	uuidString := r.PathValue("uuid")
	if err := uuid.Validate(uuidString); err != nil {
		http.Error(w, "GET /api/articles/{uuid} error: "+uuidString, http.StatusBadRequest)
		return
	}
	article, err := h.DI.DbCon.ReadArticleByUUID(uuidString)
	if err != nil {
		http.Error(w, "GET /api/articles/{uuid} error", http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(article); err != nil {
		http.Error(w, "GET /api/articles/{uuid} error: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) GetArticleThumbnail(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("uuid")
	if err := uuid.Validate(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	b, err := h.DI.DbCon.ReadArticleThumnail(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(*b) == 0 {
		http.Error(w, "no image", http.StatusInternalServerError)
		return
	}
	w.Write(*b)
}
