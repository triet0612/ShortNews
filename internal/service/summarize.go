package service

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"newscrapper/internal/db"
	"newscrapper/internal/model"
	"os/exec"
	"regexp"
	"strings"
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

func hostOllama() {
	res, err := http.Get("http://localhost:11434")
	if err == nil && res.StatusCode == 200 {
		slog.Info("ollama started")
	} else {
		slog.Info("waiting to start ollama")
		go exec.Command("ollama", "serve").Run()
		for {
			res, err := http.Get("http://localhost:11434")
			if err != nil {
				time.Sleep(5 * time.Second)
				continue
			}
			if res.StatusCode == 200 {
				break
			}
			time.Sleep(5 * time.Second)
		}
	}
}

func NewSummarizeService(db *db.DBService, audio *AudioService, config map[string]string) *SummarizeService {
	hostOllama()
	llm, err := ollama.New(ollama.WithModel("gemma:2b-instruct-v1.1-q4_0"))
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
				continue
			}
		}
	}
}

func (s *SummarizeService) ReadArticleNoSummary() (*[]model.Article, error) {
	ans := []model.Article{}
	rows, err := s.db.QueryContext(
		context.Background(),
		`SELECT a.*, n.Language FROM Article a JOIN NewsSource n
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
		a.Ext["Language"] = temp
		ans = append(ans, *a)
	}
	if len(ans) == 0 {
		return nil, errors.New("no items")
	}
	return &ans, nil
}

func (s *SummarizeService) llmSummarize(doc string, language string) (string, error) {
	doc = cleanText(doc)
	prompt := `Can you provide a comprehensive summary of the given text? The summary should cover all the key points and main ideas presented in the original text, while also condensing the information into a concise and easy-to-understand format. Please ensure that the summary includes relevant details and examples that support the main ideas, while avoiding any unnecessary information or repetition. The length of the summary should be appropriate for the length and complexity of the original text, providing a clear and accurate overview without omitting any important information.`
	if language == "vi" {
		prompt = "Bạn có thể  tóm tắt văn bản trên? Tóm tắt nên bao gồm những thông tin chính trong văn bản gốc và cô đọng những thông tin đó một cách dễ hiểu và ngắn gọn. Hãy đảm bảo những chi tiết quan trọng theo ý của tác giả và tránh những thông tin không cần thiết hay lặp lại. Độ dài nên phù hợp với độ dài và độ phức tạp của văn bản gốc, giữ các thông tin chính xác và ngắn gọn mà không loại bỏ những thông tin quan trọng."
	}
	doc = doc + "\n\n" + prompt

	ctx := context.Background()

	textResponse, err := s.llm.Call(ctx, doc,
		llms.WithTemperature(0), llms.WithTopP(1), llms.WithMaxTokens(512),
		llms.WithMaxLength(512), llms.WithFrequencyPenalty(0), llms.WithPresencePenalty(0),
	)
	if err != nil {
		return "", err
	}
	textResponse = cleanText(textResponse)
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

func cleanText(doc string) string {
	re := regexp.MustCompile(`\*\*.+\*\*`)
	doc = re.ReplaceAllString(doc, "")
	re = regexp.MustCompile(`[\[\];'":<>/|\\=+\-_()*&^%$#@!~]`)
	doc = re.ReplaceAllString(doc, " ")
	doc = strings.Trim(doc, " ")
	return doc
}
