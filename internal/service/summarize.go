package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"newscrapper/internal/db"
	"newscrapper/internal/model"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

type SummarizeService struct {
	db    *db.DBService
	llm   *ollama.LLM
	audio *AudioService
}

func NewSummarizeService(db *db.DBService, llm *ollama.LLM, audio *AudioService) *SummarizeService {
	return &SummarizeService{db: db, llm: llm, audio: audio}
}

func (s *SummarizeService) ArticleSummarize(ctx context.Context) {
	articles, err := s.ReadArticleNoSummary()
	if err != nil {
		slog.Warn(err.Error())
		return
	}
	client := http.Client{Timeout: 10 * time.Second}
	for _, article := range *articles {
		select {
		case <-ctx.Done():
			slog.Info("Ending Article Summary")
			return
		default:
			res, err := client.Get(article.Link)
			if err != nil {
				slog.Warn(err.Error())
				continue
			}
			doc, err := goquery.NewDocumentFromReader(res.Body)
			if err != nil {
				slog.Error(err.Error())
				continue
			}
			article.Summary = ""
			doc.Find("p").Each(func(i int, s *goquery.Selection) {
				if i == 0 {
					return
				}
				article.Summary += s.Text()
			})
			lang := article.Ext["Language"].(string)
			if article.Summary, err = s.llmSummarize(article.Summary, lang); err != nil {
				slog.Warn(err.Error())
				continue
			}
			if err := s.updateSummaryArticle(&article); err != nil {
				slog.Warn(err.Error())
				continue
			}
			if err := s.audio.updateArticleAudio(article.ArticleID, article.Summary, lang, article.Title); err != nil {
				slog.Warn(err.Error())
				return
			}
		}
	}
}

func (s *SummarizeService) ReadArticleNoSummary() (*[]model.Article, error) {
	ans := []model.Article{}
	rows, err := s.db.QueryContext(
		context.Background(),
		`SELECT a.*, n.Language FROM Article a JOIN NewsSource n
ON a.PublisherID=n.PublisherID WHERE a.Summary = ''`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		a := &model.Article{}
		a.Ext = map[string]interface{}{}
		temp := ""
		if err := rows.Scan(&a.ArticleID, &a.Link, &a.Title, &a.PubDate, &a.PublisherID, &a.Summary, &temp); err != nil {
			slog.Warn(err.Error())
			continue
		}
		a.Ext["Language"] = temp
		ans = append(ans, *a)
	}
	if len(ans) == 0 {
		return nil, errors.New("no items")
	}
	return &ans, nil
}

func (s *SummarizeService) llmSummarize(doc string, language string) (string, error) {
	doc = doc + "\n" + fmt.Sprintf(`Explain the above in one paragraph with %s language.`, language)

	ctx := context.Background()

	textResponse, err := s.llm.Call(ctx, doc,
		llms.WithTemperature(1), llms.WithTopP(1), llms.WithMaxTokens(250),
		llms.WithMaxLength(600), llms.WithFrequencyPenalty(0), llms.WithPresencePenalty(0),
	)

	if err != nil {
		return "", err
	}
	return textResponse, nil
}

func (s *SummarizeService) updateSummaryArticle(article *model.Article) error {
	if _, err := s.db.ExecContext(context.Background(),
		"UPDATE Article SET Summary = ? WHERE ArticleID = ?",
		article.Summary, article.ArticleID,
	); err != nil {
		return err
	}
	return nil
}
