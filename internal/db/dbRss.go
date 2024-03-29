package db

import (
	"context"
	"fmt"
	"newscrapper/model"
)

func (d *DBService) ReadSourceRSS() (*[]model.NewsSource, error) {
	ans := []model.NewsSource{}

	rows, err := d.QueryContext(context.Background(), "SELECT * FROM NewsSource")
	if err != nil {
		return nil, fmt.Errorf("ReadSourceRSS error: %s", err)
	}

	for rows.Next() {
		var temp model.NewsSource
		err := rows.Scan(&temp.PublisherID, &temp.Publisher, &temp.Link, &temp.Language)
		if err != nil {
			return nil, fmt.Errorf("ReadSourceRSS error: %s", err)
		}
		ans = append(ans, temp)
	}

	if len(ans) == 0 {
		return nil, fmt.Errorf("ReadSourceRSS error: no source found")
	}
	return &ans, nil
}

func (d *DBService) InsertRSSsource(src model.NewsSource) error {
	if _, err := d.ExecContext(context.Background(),
		"INSERT INTO NewsSource VALUES (?, ?, ?, ?)",
		src.PublisherID, src.Publisher, src.Link, src.Language,
	); err != nil {
		return fmt.Errorf("InsertArticle error: %s", err)
	}
	return nil
}

func (d *DBService) UpdateRSSsource(src model.NewsSource) error {
	if _, err := d.ExecContext(context.Background(),
		"UPDATE NewsSource SET Link=?, Language=?, Publisher=? WHERE UUID=?",
		src.Link, src.Language, src.Publisher, src.PublisherID,
	); err != nil {
		return fmt.Errorf("UpdateRSSsource error: %s", err)
	}
	return nil
}

func (d *DBService) DeleteRSSsource(src model.NewsSource) error {
	if _, err := d.ExecContext(context.Background(),
		"DELETE FROM NewsSource WHERE UUID=?", src.PublisherID,
	); err != nil {
		return fmt.Errorf("DeleteRSSsource error: %s", err)
	}
	return nil
}
