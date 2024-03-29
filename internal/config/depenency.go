package config

import (
	"log"
	"newscrapper/internal/db"
	"time"

	"github.com/tmc/langchaingo/llms/ollama"
)

type DI struct {
	DbCon   *db.DBService
	Llmodel *ollama.LLM
	Clock   *Clock
}

func InitDependency() *DI {
	dbService := db.GetInstance()
	llm, err := ollama.New(ollama.WithModel("llama2"))
	if err != nil {
		log.Fatal(err)
	}
	return &DI{
		DbCon:   dbService,
		Llmodel: llm,
		Clock: &Clock{
			Timer:    time.NewTimer(1 * time.Second),
			PollRate: 1 * time.Hour,
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
