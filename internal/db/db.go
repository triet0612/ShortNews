package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type DBService struct {
	*sql.Conn
}

const SqlScript = `
CREATE TABLE Article (
	ArticleID TEXT PRIMARY KEY NOT NULL,
	Link TEXT UNIQUE NOT NULL,
	Title TEXT NOT NULL,
	PubDate TEXT NOT NULL,
	Publisher TEXT NOT NULL,
	Summary TEXT NOT NULL DEFAULT ('')
);
CREATE TABLE Thumbnail (
	ArticleID TEXT PRIMARY KEY NOT NULL,
	URL TEXT NOT NULL,
	Image BLOB NOT NULL
);
CREATE TABLE NewsSource (
	PublisherID TEXT PRIMARY KEY NOT NULL,
	Publisher TEXT UNIQUE NOT NULL,
	Link TEXT NOT NULL,
	Language TEXT NOT NULL
);
CREATE TABLE Config (
	KEY TEXT PRIMARY KEY,
	VALUE TEXT NOT NULL
);
CREATE INDEX "INDEX_ARTICLE" ON "Article" (
	"PubDate" DESC
);
INSERT INTO Config VALUES ("RSS-refresh-rate", "5");`

func GetInstance() *DBService {
	_, err := os.Stat("./news.db")
	not_exist := errors.Is(err, os.ErrNotExist)
	db, err := sql.Open("sqlite3", "./news.db")
	if err != nil {
		log.Panic(fmt.Errorf("GetInstance error: %s", err))
	}
	if not_exist {
		if _, err := db.Exec(SqlScript); err != nil {
			log.Panic(fmt.Errorf("GetInstance error: %s", err))
		}
	}
	con, err := db.Conn(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	return &DBService{Conn: con}
}
