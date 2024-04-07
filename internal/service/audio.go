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
	db        *db.DBService
	langAudio map[string]string
}

func NewAudioService(db *db.DBService, langAudio map[string]string) *AudioService {
	return &AudioService{db: db, langAudio: langAudio}
}

func (a *AudioService) GenerateAudio() {
	articles, err := a.readArticleNoAudio()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	for _, article := range *articles {
		go func() {
			id, sum, lang, title := article[0], article[1], article[2], article[3]
			if err = a.updateArticleAudio(id, sum, lang, title); err != nil {
				slog.Error(err.Error())
				return
			}
		}()
	}
}

func (a *AudioService) updateArticleAudio(id string, sum string, lang string, title string) error {
	cmd := exec.Command("mimic3", "--voice", a.langAudio[lang])
	stdin, err := cmd.StdinPipe()
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	_, err = io.WriteString(stdin, title+"\n"+sum)
	if err != nil {
		slog.Warn(err.Error())
		return err
	}
	stdin.Close()
	cmd.Wait()
	out, err := cmd.Output()
	if err != nil {
		slog.Warn(err.Error())
		return err
	}
	if len(out) == 0 {
		return nil
	}
	if _, err := a.db.ExecContext(context.Background(),
		"UPDATE ArticleAudio SET Audio = ? WHERE ArticleID = ?",
		out, id,
	); err != nil {
		return err
	}
	return nil
}

func (a *AudioService) readArticleNoAudio() (*[][]string, error) {
	ans := [][]string{}
	rows, err := a.db.QueryContext(context.Background(),
		`SELECT a1.ArticleID, a1.Summary, n.Language, a1.Title
FROM Article a1 JOIN ArticleAudio a2 JOIN NewsSource n
ON a1.ArticleID=a2.ArticleID AND a1.PublisherID=n.PublisherID AND a2.Audio=""`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		id, sum, lang, title := "", "", "", ""
		if err := rows.Scan(&id, &sum, &lang, &title); err != nil {
			slog.Warn(err.Error())
			continue
		}
		ans = append(ans, []string{id, sum, lang, title})
	}
	if len(ans) == 0 {
		return nil, errors.New("no items")
	}
	return &ans, nil
}
