package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

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

func (h *Handler) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	// TODO: Update article by ID
}

func (h *Handler) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	// TODO: Delete article by ID, also delete thumbnail
}
