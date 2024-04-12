package service

import (
	"context"
	"log/slog"
	"newscrapper/internal/db"
)

type DBCleanUp struct {
	db *db.DBService
}

func NewDBCleanUp(db *db.DBService) *DBCleanUp {
	return &DBCleanUp{db: db}
}

func (s DBCleanUp) CleanOldArticle(ctx context.Context) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		slog.Error(err.Error())
		slog.Error(tx.Rollback().Error())
		return
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM ArticleAudio WHERE
ArticleID IN (
	SELECT ArticleID FROM Article WHERE date("now") - date(PubDate) > 1
)`); err != nil {
		slog.Error(err.Error())
		slog.Error(tx.Rollback().Error())
		return
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM Thumbnail WHERE
ArticleID IN (
	SELECT ArticleID FROM Article WHERE date("now") - date(PubDate) > 1
)`); err != nil {
		slog.Error(err.Error())
		slog.Error(tx.Rollback().Error())
		return
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM Article WHERE date("now") - date(PubDate) > 1`); err != nil {
		slog.Error(err.Error())
		slog.Error(tx.Rollback().Error())
		return
	}
	if err := tx.Commit(); err != nil {
		slog.Error(err.Error())
	}
}
