package service

import (
	"context"
	"errors"
	"log/slog"
	"newscrapper/internal/db"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type AudioService struct {
	db        *db.DBService
	langAudio map[string]string
	config    map[string]string
}

func NewAudioService(db *db.DBService, config map[string]string) *AudioService {
	vm := map[string]string{}
	rows, _ := db.QueryContext(context.Background(), "SELECT * FROM VoiceModel")
	for rows.Next() {
		lang, modelName := "", ""
		rows.Scan(&lang, &modelName)
		vm[lang] = modelName
	}
	return &AudioService{db: db, config: config, langAudio: vm}
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
				continue
			}
		}
	}
}

func (a *AudioService) updateArticleAudio(id string, sum string, lang string, title string) error {
	slog.Info(a.langAudio[lang])
	cmd := exec.Command(
		"piper", "--model", a.langAudio[lang],
		"--output_file", "./temp.wav")
	cmd.Stdin = strings.NewReader(cleanTextAudio(title + " " + sum))
	s := &strings.Builder{}
	cmd.Stderr = s
	if err := cmd.Run(); err != nil {
		slog.Info(cleanTextAudio(title + " " + sum))
		slog.Info(s.String())
		return err
	}
	file, err := os.ReadFile("./temp.wav")
	if err != nil {
		return err
	}

	if _, err := a.db.ExecContext(context.Background(),
		"UPDATE ArticleAudio SET Audio=? WHERE ArticleID=?",
		string(file), id,
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
