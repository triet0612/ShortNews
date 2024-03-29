package db

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"newscrapper/model"
	"strings"

	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
)

func (d *DBService) InsertArticle(a *model.RSSItem) error {
	a.ArticleID = uuid.New().String()
	trimdate := strings.Split(a.PubDate, " ")
	a.PubDate = strings.Join(trimdate[1:len(trimdate)-1], " ")

	ctx := context.Background()

	if _, err := d.ExecContext(ctx,
		`INSERT INTO Article (ArticleID, Link, Title, PubDate, Publisher)
		VALUES (?, ?, ?, ?, ?)`,
		a.ArticleID, a.Link, a.Title, a.PubDate, a.Publisher,
	); err != nil {
		var e sqlite3.Error
		if ok := errors.As(err, &e); ok && (e.ExtendedCode == sqlite3.ErrConstraintUnique) {
			return nil
		}
		return err
	}
	if a.Image.URL == "" {
		return nil
	}
	if _, err := d.ExecContext(ctx,
		"INSERT INTO Thumbnail VALUES (?, ?, ?)",
		a.ArticleID, a.Image.URL, "",
	); err != nil {
		return err
	}
	return nil
}

func (d *DBService) UpdateSummaryArticle(article *model.RSSItem) error {
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

func (d *DBService) ReadArticle(limit uint, offset uint) (*[]model.RSSItem, error) {
	ans := []model.RSSItem{}

	rows, err := d.QueryContext(context.Background(),
		"SELECT * FROM Article ORDER BY PubDate LIMIT ? OFFSET ?",
		limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		a := &model.RSSItem{}
		if err := rows.Scan(&a.ArticleID, &a.Link, &a.Title, &a.PubDate, &a.Publisher, &a.Summary); err != nil {
			slog.Warn(err.Error())
			continue
		}
		ans = append(ans, *a)
	}
	return &ans, nil
}

func (d *DBService) ReadArticleNoSummary() (*[]model.RSSItem, error) {
	ans := []model.RSSItem{}
	rows, err := d.QueryContext(context.Background(), "SELECT * FROM Article WHERE Summary = ''")
	if err != nil {
		return nil, fmt.Errorf("ReadArticleNoSummary: %s", err)
	}
	defer rows.Close()
	for rows.Next() {
		a := &model.RSSItem{}
		if err := rows.Scan(&a.ArticleID, &a.Link, &a.Title, &a.PubDate, &a.Publisher, &a.Summary); err != nil {
			slog.Warn(err.Error())
			continue
		}
		ans = append(ans, *a)
	}
	if len(ans) == 0 {
		return nil, errors.New("ReadArticleNoSummary: no items")
	}
	return &ans, nil
}

func (d *DBService) ReadArticleByUUID(id string) (*model.RSSItem, error) {
	a := &model.RSSItem{}
	row := d.QueryRowContext(context.Background(), "SELECT * FROM Article WHERE UUID=?", id)
	if err := row.Scan(&a.ArticleID, &a.Link, &a.Title, &a.PubDate, &a.Publisher, &a.Summary); err != nil {
		return nil, err
	}
	return a, nil
}

func (d *DBService) ReadArticleNoThumbnail() (*[][]string, error) {
	ans := [][]string{}

	rows, err := d.QueryContext(context.Background(),
		"SELECT ArticleID, URL FROM Thumbnail WHERE URL IS NOT NULL AND Image = ''",
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
	row := d.QueryRowContext(context.Background(), "SELECT Thumbnail FROM Article WHERE UUID=?", id)
	if err := row.Scan(&ans); err != nil {
		return nil, err
	}
	return &ans, nil
}
