package db

import (
	"context"
	"errors"
	"newscrapper/internal/model"
)

func (d *DBService) ReadSourceRSS() (*[]model.NewsSource, error) {
	ans := []model.NewsSource{}
	rows, err := d.QueryContext(context.Background(), "SELECT * FROM NewsSource")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var temp model.NewsSource
		err := rows.Scan(&temp.PublisherID, &temp.Link, &temp.Language)
		if err != nil {
			return nil, err
		}
		ans = append(ans, temp)
	}
	if len(ans) == 0 {
		return nil, errors.New("no source found")
	}
	return &ans, nil
}
