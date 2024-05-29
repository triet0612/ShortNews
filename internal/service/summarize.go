package service

import (
	"context"
	"errors"
	"log"
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

func NewSummarizeService(db *db.DBService, audio *AudioService, config map[string]string) *SummarizeService {
	llm, err := ollama.New(ollama.WithServerURL(config["ollama_api"]), ollama.WithModel("gemma:2b-instruct-v1.1-q4_0"))
	if err != nil {
		log.Fatal(err)
	}
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
			if article.Summary, err = s.llmSummarize(article.Summary); err != nil {
				slog.Warn(err.Error())
				continue
			}
			if err := s.updateSummaryArticle(&article); err != nil {
				slog.Warn(err.Error())
				continue
			}
			if err := s.audio.updateArticleAudio(article.ArticleID, article.Summary, article.Title, article.Ext["VoiceType"].(string)); err != nil {
				slog.Warn(err.Error())
				continue
			}
		}
	}
}

func (s *SummarizeService) ReadArticleNoSummary() (*[]model.Article, error) {
	ans := []model.Article{}
	rows, err := s.db.QueryContext(
		context.Background(),
		`SELECT a.*, n.VoiceType FROM Article a JOIN NewsSource n
	ON a.PublisherID=n.PublisherID WHERE a.Summary = '' ORDER BY a.PubDate DESC`,
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
		a.Ext["VoiceType"] = temp
		ans = append(ans, *a)
	}
	if len(ans) == 0 {
		return nil, errors.New("no items")
	}
	return &ans, nil
}

func (s *SummarizeService) llmSummarize(doc string) (string, error) {
	prompt := "Bạn có thể  tóm tắt văn bản trên? Tóm tắt nên bao gồm những thông tin chính trong văn bản gốc và cô đọng những thông tin đó một cách dễ hiểu và ngắn gọn. Hãy đảm bảo những chi tiết quan trọng theo ý của tác giả và tránh những thông tin không cần thiết hay lặp lại. Độ dài nên phù hợp với độ dài và độ phức tạp của văn bản gốc, giữ các thông tin chính xác và ngắn gọn mà không loại bỏ những thông tin quan trọng."
	doc = doc + "\n\n" + prompt

	ctx := context.Background()

	textResponse, err := s.llm.Call(ctx, doc,
		llms.WithTemperature(0), llms.WithTopP(1), llms.WithMaxTokens(256),
		llms.WithMaxLength(256), llms.WithFrequencyPenalty(0), llms.WithPresencePenalty(0),
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
