package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"newscrapper/internal/db"
	"os/exec"
	"strings"
	"time"
)

type AudioService struct {
	db        *db.DBService
	langAudio map[string]string
}

func NewAudioService(db *db.DBService) *AudioService {
	hostMimic3()
	vm := map[string]string{}
	rows, _ := db.QueryContext(context.Background(), "SELECT * FROM VoiceModel")
	for rows.Next() {
		lang, modelName := "", ""
		rows.Scan(&lang, &modelName)
		vm[lang] = modelName
	}
	return &AudioService{db: db, langAudio: vm}
}

func hostMimic3() {
	slog.Info("starting mimic3-server")
	go exec.Command("mimic3-server", "--port", "8001").Run()
	for {
		res, err := http.Get("http://localhost:8001")
		if err != nil {
			time.Sleep(5 * time.Second)
			continue
		}
		if res.StatusCode == 200 {
			break
		}
		time.Sleep(10 * time.Second)
	}
	slog.Info("starting mimic3-server complete")
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
			id, sum, lang, title := article[0], article[1], article[2], article[3]
			sum = cleanTextAudio(sum)
			if err = a.updateArticleAudio(id, sum, lang, title); err != nil {
				slog.Error(err.Error())
				return
			}
		}
	}
}

func (a *AudioService) updateArticleAudio(id string, sum string, lang string, title string) error {
	client := http.Client{Timeout: 10 * time.Second}
	res, err := client.Post(fmt.Sprintf("http://localhost:8001/api/tts?voice=%s", a.langAudio[lang]), "text", strings.NewReader(title+"\n"+sum))
	if err != nil {
		return err
	}
	out, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if len(out) == 0 {
		return nil
	}
	if _, err := a.db.ExecContext(context.Background(),
		"UPDATE ArticleAudio SET Audio=? WHERE ArticleID=?",
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

func cleanTextAudio(doc string) string {
	doc = strings.ReplaceAll(doc, "\n", "        ")
	doc = strings.ReplaceAll(doc, ".", "    ")
	doc = strings.ReplaceAll(doc, ",", "  ")
	return doc
}
