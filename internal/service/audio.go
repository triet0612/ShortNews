package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"log/slog"
	"net/http"
	"newscrapper/internal/db"
)

type AudioService struct {
	db     *db.DBService
	config map[string]string
}

func NewAudioService(db *db.DBService, config map[string]string) *AudioService {
	return &AudioService{db: db, config: config}
}

func (a *AudioService) GenerateAudio(ctx context.Context) {
	articles, err := a.readArticleNoAudio()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	for _, article := range *articles {
		select {
		case <-ctx.Done():
			slog.Info("Ending GenerateAudio")
			return
		default:
			id, sum, voice, title := article[0], article[1], article[2], article[3]
			if err = a.updateArticleAudio(id, sum, title, voice); err != nil {
				slog.Error(err.Error())
				continue
			}
		}
	}
}

func (a *AudioService) updateArticleAudio(id string, sum string, title string, voice string) error {
	body := map[string]string{
		"text": title + "\n" + sum,
		"type": voice,
	}
	b, _ := json.Marshal(body)
	res, err := http.Post(a.config["voice_api"]+"/text-to-speech/",
		"application/json",
		bytes.NewBuffer(b),
	)
	log.Println(body)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	audio, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if len(audio) == 0 {
		return nil
	}
	if _, err := a.db.ExecContext(context.Background(),
		"UPDATE ArticleAudio SET Audio=? WHERE ArticleID=?",
		string(audio), id,
	); err != nil {
		return err
	}
	return nil
}

func (a *AudioService) readArticleNoAudio() (*[][]string, error) {
	ans := [][]string{}
	rows, err := a.db.QueryContext(context.Background(),
		`SELECT a1.ArticleID, a1.Summary, n.VoiceType, a1.Title
FROM Article a1 JOIN ArticleAudio a2 JOIN NewsSource n
ON a1.ArticleID=a2.ArticleID AND a1.PublisherID=n.PublisherID AND a2.Audio=""`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		id, sum, voice, title := "", "", "", ""
		if err := rows.Scan(&id, &sum, &voice, &title); err != nil {
			slog.Warn(err.Error())
			continue
		}
		ans = append(ans, []string{id, sum, voice, title})
	}
	if len(ans) == 0 {
		return nil, errors.New("no items")
	}
	return &ans, nil
}
