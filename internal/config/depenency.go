package config

import (
	"context"
	"log"
	"newscrapper/internal/db"
	"time"

	"github.com/tmc/langchaingo/llms/ollama"
)

type DI struct {
	DbCon     *db.DBService
	Llmodel   *ollama.LLM
	Clock     *Clock
	LangAudio map[string]string
}

func InitDependency() *DI {
	dbService := db.GetInstance()
	llm, err := ollama.New(ollama.WithModel("gemma:2b"))
	if err != nil {
		log.Fatal(err)
	}
	vm := map[string]string{}
	rows, _ := dbService.QueryContext(context.Background(), "SELECT * FROM VoiceModel")
	for rows.Next() {
		lang, modelName := "", ""
		rows.Scan(&lang, &modelName)
		vm[lang] = modelName
	}

	return &DI{
		DbCon:   dbService,
		Llmodel: llm,
		Clock: &Clock{
			Timer:    time.NewTimer(1 * time.Second),
			PollRate: 1 * time.Hour,
		},
		LangAudio: vm,
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
