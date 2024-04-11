package service

import (
	"context"
	"errors"
	"log/slog"
	"newscrapper/internal/db"
	"newscrapper/internal/model"
	"runtime"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
	"github.com/mmcdole/gofeed"
)

type RSSFetchService struct {
	db *db.DBService
}

func NewRSSFetchService(db *db.DBService) *RSSFetchService {
	return &RSSFetchService{db: db}
}

func (r *RSSFetchService) NewsExtraction(ctx context.Context) {
	sourceList, err := r.db.ReadSourceRSS()
	if err != nil {
		slog.Warn(err.Error())
		return
	}
	wg := sync.WaitGroup{}
	articlepool := make(chan struct{}, runtime.NumCPU())
	for _, source := range *sourceList {
		select {
		case <-ctx.Done():
			return
		default:
			ctxTimeOut, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			fp := gofeed.NewParser()
			feed, err := fp.ParseURLWithContext(source.Link, ctxTimeOut)
			if err != nil {
				slog.Error(err.Error())
				continue
			}

			articlepool <- struct{}{}
			wg.Add(1)
			go func() {
				defer wg.Done()
				for _, article := range feed.Items {

					a := model.Article{
						ArticleID:   uuid.NewString(),
						Link:        article.Link,
						Title:       article.Title,
						PubDate:     article.PublishedParsed,
						PublisherID: source.PublisherID,
					}
					if article.PublishedParsed == nil {
						t := time.Now()
						a.PubDate = &t
					}
					imageLink := ""
					if article.Image != nil {
						imageLink = (*article.Image).URL
					} else {
						temp, ok := article.Extensions["media"]["thumbnail"]
						if ok {
							imageLink = temp[0].Attrs["url"]
						}
					}
					if err := r.insertArticle(&a, imageLink); err != nil {
						slog.Info(err.Error())
						continue
					}
				}
				<-articlepool
			}()
		}
	}
	wg.Wait()
}

func (r *RSSFetchService) insertArticle(a *model.Article, imageUrl string) error {
	ctx := context.Background()
	for _, err := r.db.ExecContext(ctx,
		`INSERT INTO Article (ArticleID, Link, Title, PubDate, PublisherID) VALUES (?, ?, ?, ?, ?)`,
		a.ArticleID, a.Link, a.Title, a.PubDate.Format(time.DateTime), a.PublisherID,
	); err != nil; a.ArticleID = uuid.NewString() {
		var e sqlite3.Error
		if ok := errors.As(err, &e); ok {
			if e.ExtendedCode == sqlite3.ErrConstraintUnique {
				return nil
			} else if e.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
				continue
			}
		}
		return err
	}
	if imageUrl == "" {
		return nil
	}
	if _, err := r.db.ExecContext(ctx,
		"INSERT INTO Thumbnail VALUES (?, ?)",
		a.ArticleID, imageUrl,
	); err != nil {
		return err
	}
	if _, err := r.db.ExecContext(ctx, "INSERT INTO ArticleAudio VALUES (?, ?)", a.ArticleID, ""); err != nil {
		return err
	}
	return nil
}
