package service

import (
	"context"
	"encoding/xml"
	"io"
	"log"
	"log/slog"
	"net/http"
	"newscrapper/internal/config"
	"newscrapper/model"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

type RssPullService struct {
	Di *config.DI
}

func (r *RssPullService) RunRssService() {
	for {
		<-r.Di.Clock.Timer.C
		r.Di.Clock.Timeout()
		r.NewsExtraction()
		go r.UpdateThumbnail()
		go r.ArticleSummarize()
	}
}

func (r *RssPullService) NewsExtraction() {
	sourceList, err := r.Di.DbCon.ReadSourceRSS()
	if err != nil {
		slog.Warn(err.Error())
		return
	}
	wg := sync.WaitGroup{}
	articlepool := make(chan struct{}, runtime.NumCPU())
	for _, source := range *sourceList {
		srcClient := http.Client{Timeout: 10 * time.Second}
		res, err := srcClient.Get(source.Link)
		if err != nil {
			slog.Error(err.Error())
			continue
		}

		rss := model.RSSNews{}
		if err := xml.NewDecoder(res.Body).Decode(&rss); err != nil {
			log.Println("NewsExtraction error: ", err)
			continue
		}

		articlepool <- struct{}{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, article := range rss.Channel.Items {
				article.Publisher = source.Publisher

				if err := r.Di.DbCon.InsertArticle(&article); err != nil {
					slog.Info(err.Error())
					continue
				}
			}
			<-articlepool
		}()
	}
	wg.Wait()
}

func (r *RssPullService) ArticleSummarize() {
	articles, err := r.Di.DbCon.ReadArticleNoSummary()
	if err != nil {
		slog.Warn(err.Error())
		return
	}

	client := http.Client{Timeout: 10 * time.Second}
	for _, article := range *articles {
		res, err := client.Get(article.Link)
		if err != nil {
			slog.Warn(err.Error())
			return
		}
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			slog.Error("ArticleSummarize error: %s", err)
			continue
		}
		article.Summary = ""
		doc.Find("p").Each(func(i int, s *goquery.Selection) {
			article.Summary += s.Text()
		})

		if len(strings.Split(article.Summary, " ")) > 250 {
			if article.Summary, err = llmSummarize(
				r.Di.Llmodel, article.Summary,
			); err != nil {
				slog.Warn(err.Error())
				return
			}
		}

		if err := r.Di.DbCon.UpdateSummaryArticle(&article); err != nil {
			slog.Warn(err.Error())
			return
		}
	}
}

func (r *RssPullService) UpdateThumbnail() {
	articles, err := r.Di.DbCon.ReadArticleNoThumbnail()
	if err != nil {
		slog.Warn(err.Error())
		return
	}
	client := http.Client{Timeout: 10 * time.Second}
	for _, a := range *articles {
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
		if err := r.Di.DbCon.UpdateArticleThumbnail(id, img); err != nil {
			slog.Warn(err.Error())
		}
	}
}

func llmSummarize(llm *ollama.LLM, doc string) (string, error) {
	doc = doc + "\n" + `Return main content of the document above using the same language as the document with less than two-hundred words.`

	ctx := context.Background()

	textResponse, err := llm.Call(ctx, doc, llms.WithTemperature(0))

	if err != nil {
		return "", err
	}
	re, err := regexp.Compile("[~!@#$%^&*()\\-_+={}\\]\\[\\|\\`,./?;:'\"<>^]{2,}")
	if err != nil {
		return "", err
	}
	ans := re.ReplaceAllString(textResponse, "")
	return ans, nil
}
