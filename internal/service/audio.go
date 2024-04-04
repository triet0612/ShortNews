package service

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"newscrapper/internal/db"
	"os/exec"
)

type AudioService struct {
	db *db.DBService
}

func NewAudioService(db *db.DBService) *AudioService {
	return &AudioService{db: db}
}

func (a *AudioService) GenerateAudio() {
	articles, err := a.readArticleNoAudio()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	for _, article := range *articles {
		id, sum, lang := article[0], article[1], article[2]

		cmd := exec.Command("mimic3", "--voice", lang)

		stdin, err := cmd.StdinPipe()
		if err != nil {
			continue
		}
		_, err = io.WriteString(stdin, sum)
		if err != nil {
			slog.Warn(err.Error())
		}
		stdin.Close()
		cmd.Wait()
		out, err := cmd.Output()
		if err != nil {
			slog.Warn(err.Error())
			continue
		}
		if err = a.updateArticleAudio(id, out); err != nil {
			slog.Error(err.Error())
			continue
		}
	}
}

func (a *AudioService) updateArticleAudio(id string, audio []byte) error {
	if len(audio) == 0 {
		return nil
	}
	if _, err := a.db.ExecContext(context.Background(),
		"UPDATE ArticleAudio SET Audio = ? WHERE ArticleID = ?",
		audio, id,
	); err != nil {
		return err
	}
	return nil
}

func (a *AudioService) readArticleNoAudio() (*[][]string, error) {
	ans := [][]string{}
	rows, err := a.db.QueryContext(context.Background(),
		`SELECT a1.ArticleID, a1.Summary, v.ModelName
FROM Article a1 JOIN ArticleAudio a2 JOIN NewsSource n JOIN VoiceModel v
ON v.Language=n.Language AND a1.ArticleID=a2.ArticleID AND a1.PublisherID = n.PublisherID`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		id, sum, lang := "", "", ""
		if err := rows.Scan(&id, &sum, &lang); err != nil {
			slog.Warn(err.Error())
			continue
		}
		ans = append(ans, []string{id, sum, lang})
	}
	if len(ans) == 0 {
		return nil, errors.New("no items")
	}
	return &ans, nil
}
