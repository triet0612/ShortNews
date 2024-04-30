package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"newscrapper/internal/db"
	"regexp"
	"strings"
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
			id, sum := article[0], article[1]
			sum = cleanTextAudio(sum)
			if err = a.updateArticleAudio(id, sum); err != nil {
				slog.Error(err.Error())
				continue
			}
		}
	}
}

func (a *AudioService) updateArticleAudio(id string, sum string) error {
	body := map[string]string{
		"text": cleanTextAudio(sum),
	}
	b, _ := json.Marshal(body)
	res, err := http.Post(a.config["voice_api"]+"/text-to-speech/",
		"application/json",
		bytes.NewBuffer(b),
	)
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
	reg := regexp.MustCompile(`[^a-z0-9A-Z\s_ÀÁÂÃÈÉÊÌÍÒÓÔÕÙÚĂĐĨŨƠàáâãèéêìíòóôõùúăđĩũơƯĂẠẢẤẦẨẪẬẮẰẲẴẶẸẺẼỀỀỂưăạảấầẩẫậắằẳẵặẹẻẽềềểỄỆỈỊỌỎỐỒỔỖỘỚỜỞỠỢỤỦỨỪễếệỉịọỏốồổỗộớờởỡợụủứừỬỮỰỲỴÝỶỸửữựỳỵỷỹ]`)
	doc = reg.ReplaceAllString(doc, "")
	return strings.ReplaceAll(doc, "\n", "")
}
