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
	PubDate DATETIME NOT NULL,
	PublisherID TEXT NOT NULL,
	Summary TEXT NOT NULL DEFAULT ('')
);
CREATE TABLE Thumbnail (
	ArticleID TEXT PRIMARY KEY NOT NULL,
	URL TEXT NOT NULL,
	Image BLOB NOT NULL
);
CREATE TABLE NewsSource (
	PublisherID TEXT PRIMARY KEY NOT NULL,
	Link TEXT NOT NULL,
	Language TEXT NOT NULL
);
CREATE TABLE Config (
	KEY TEXT PRIMARY KEY,
	VALUE TEXT NOT NULL
);
CREATE TABLE VoiceModel (
	Language TEXT PRIMARY KEY,
	ModelName TEXT
);
CREATE TABLE ArticleAudio(
	ArticleID TEXT PRIMARY KEY NOT NULL,
	Audio BLOB NOT NULL
);
INSERT INTO VoiceModel VALUES ("vi", "vi_VN/vais1000_low");
INSERT INTO VoiceModel VALUES ("en", "en_US/cmu-arctic_low");`

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
