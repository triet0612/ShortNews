package db

import (
	"context"
	"errors"
	"log/slog"
	"newscrapper/model"
	"time"

	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
)

func (d *DBService) InsertArticle(a *model.Article, imageUrl string) error {
	ctx := context.Background()
	for _, err := d.ExecContext(ctx,
		`INSERT INTO Article (ArticleID, Link, Title, PubDate, Publisher) VALUES (?, ?, ?, ?, ?)`,
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
	if _, err := d.ExecContext(ctx,
		"INSERT INTO Thumbnail VALUES (?, ?, ?)",
		a.ArticleID, imageUrl, "",
	); err != nil {
		return err
	}
	return nil
}

func (d *DBService) UpdateSummaryArticle(article *model.Article) error {
	if _, err := d.ExecContext(context.Background(),
		"UPDATE Article SET Summary = ? WHERE ArticleID = ?",
		article.Summary, article.ArticleID,
	); err != nil {
		return err
	}
	return nil
}

func (d *DBService) UpdateArticleThumbnail(id string, img []byte) error {
	if len(img) == 0 {
		return nil
	}
	if _, err := d.ExecContext(context.Background(),
		"UPDATE Thumbnail SET Image = ? WHERE ArticleID = ?",
		img, id,
	); err != nil {
		return err
	}
	return nil
}

func (d *DBService) ReadArticle(limit uint, offset uint) (*[]model.Article, error) {
	ans := []model.Article{}

	rows, err := d.QueryContext(context.Background(),
		"SELECT * FROM Article ORDER BY PubDate LIMIT ? OFFSET ?",
		limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		a := &model.Article{}
		if err := rows.Scan(&a.ArticleID, &a.Link, &a.Title, &a.PubDate, &a.PublisherID, &a.Summary); err != nil {
			slog.Warn(err.Error())
			continue
		}
		ans = append(ans, *a)
	}
	return &ans, nil
}

func (d *DBService) ReadArticleNoSummary() (*[]model.Article, error) {
	ans := []model.Article{}
	rows, err := d.QueryContext(context.Background(), "SELECT * FROM Article WHERE Summary = ''")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		a := &model.Article{}
		if err := rows.Scan(&a.ArticleID, &a.Link, &a.Title, &a.PubDate, &a.PublisherID, &a.Summary); err != nil {
			slog.Warn(err.Error())
			continue
		}
		ans = append(ans, *a)
	}
	if len(ans) == 0 {
		return nil, errors.New("no items")
	}
	return &ans, nil
}

func (d *DBService) ReadArticleByUUID(id string) (*model.Article, error) {
	a := &model.Article{}
	row := d.QueryRowContext(context.Background(), "SELECT * FROM Article WHERE ArticleID=?", id)
	if err := row.Scan(&a.ArticleID, &a.Link, &a.Title, &a.PubDate, &a.PublisherID, &a.Summary); err != nil {
		return nil, err
	}
	return a, nil
}

func (d *DBService) ReadArticleNoThumbnail() (*[][]string, error) {
	ans := [][]string{}

	rows, err := d.QueryContext(context.Background(),
		"SELECT ArticleID, URL FROM Thumbnail WHERE Image = ''",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		id, url := "", ""
		if err := rows.Scan(&id, &url); err != nil {
			slog.Warn(err.Error())
			continue
		}
		ans = append(ans, []string{id, url})
	}
	return &ans, nil
}

func (d *DBService) ReadArticleThumnail(id string) (*[]byte, error) {
	ans := []byte{}
	row := d.QueryRowContext(context.Background(), "SELECT Image FROM Thumbnail WHERE ArticleID=?", id)
	if err := row.Scan(&ans); err != nil {
		return nil, err
	}
	return &ans, nil
}
