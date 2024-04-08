package service

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"newscrapper/internal/db"
	"time"
)

type ThumbnailService struct {
	db *db.DBService
}

func NewThumbnailService(db *db.DBService) *ThumbnailService {
	return &ThumbnailService{db: db}
}

func (t *ThumbnailService) UpdateThumbnail(ctx context.Context) {
	articles, err := t.readArticleNoThumbnail()
	if err != nil {
		slog.Warn(err.Error())
		return
	}
	client := http.Client{Timeout: 10 * time.Second}
	for _, a := range *articles {
		select {
		case <-ctx.Done():
			return
		default:
			id, imageUrl := a[0], a[1]
			res, err := client.Get(imageUrl)
			if err != nil {
				slog.Warn(err.Error())
				continue
			}
			defer res.Body.Close()
			img, err := io.ReadAll(res.Body)
			if err != nil {
				slog.Warn(err.Error())
				continue
			}
			if err := t.updateArticleThumbnail(id, img); err != nil {
				slog.Warn(err.Error())
			}
		}
	}
}

func (t *ThumbnailService) updateArticleThumbnail(id string, img []byte) error {
	if len(img) == 0 {
		return nil
	}
	if _, err := t.db.ExecContext(context.Background(),
		"UPDATE Thumbnail SET Image = ? WHERE ArticleID = ?",
		img, id,
	); err != nil {
		return err
	}
	return nil
}

func (t *ThumbnailService) readArticleNoThumbnail() (*[][]string, error) {
	ans := [][]string{}

	rows, err := t.db.QueryContext(context.Background(),
		"SELECT ArticleID, URL FROM Thumbnail WHERE Image = ''",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		id, url := "", ""
		if err := rows.Scan(&id, &url); err != nil {
			slog.Warn(err.Error())
			continue
		}
		ans = append(ans, []string{id, url})
	}
	return &ans, nil
}
