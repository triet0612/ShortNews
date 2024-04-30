package config

import (
	"newscrapper/internal/db"
	"os"
	"time"
)

type DI struct {
	DbCon  *db.DBService
	Clock  *Clock
	Signal chan struct{}
	Config map[string]string
}

func InitDependency() *DI {
	dbService := db.GetInstance()

	return &DI{
		DbCon: dbService,
		Clock: &Clock{
			Timer:    time.NewTimer(1 * time.Second),
			PollRate: 1 * time.Hour,
		},
		Signal: make(chan struct{}),
		Config: map[string]string{
			"ollama_api": os.Getenv("ollama_api"),
			"voice_api":  os.Getenv("voice_api"),
		},
	}
}

type Clock struct {
	Timer    *time.Timer
	PollRate time.Duration
}

func (c *Clock) Timeout() {
	c.Timer.Reset(c.PollRate)
}

func (c *Clock) Sync() {
	c.Timer.Reset(time.Millisecond)
}
