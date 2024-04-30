package start

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"newscrapper/internal/config"
	"newscrapper/internal/handler"
	"newscrapper/internal/service"
	"sync"
)

func RunServices(di *config.DI) {
	go runNewsService(di)
	runHTTPServer(di)
}

func runHTTPServer(di *config.DI) {
	slog.Info("pulling model")
	body := map[string]interface{}{
		"name":   "gemma:2b-instruct-v1.1-q4_0",
		"stream": false,
	}
	b, _ := json.Marshal(body)
	_, err := http.Post(di.Config["ollama_api"]+"/api/pull", "application/json", bytes.NewBuffer(b))
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("complete pulling")
	slog.Info("start http server")
	h := handler.NewHandler(di.DbCon, di.Clock, di.Signal)
	mux := http.NewServeMux()
	h.Mount(mux)
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Fatal(server.ListenAndServe())
}

func runNewsService(di *config.DI) {
	fetcher := service.NewRSSFetchService(di.DbCon)
	audio := service.NewAudioService(di.DbCon, di.Config)
	summary := service.NewSummarizeService(di.DbCon, audio, di.Config)
	cleanup := service.NewDBCleanUp(di.DbCon)

	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	for {
		select {
		case <-di.Clock.Timer.C:
			di.Clock.Timeout()
			wg.Add(1)
			go func() {
				defer wg.Done()
				go cleanup.CleanOldArticle(ctx)
				fetcher.NewsExtraction(ctx)
				summary.ArticleSummarize(ctx)
				go audio.GenerateAudio(ctx)
			}()
		case <-di.Signal:
			cancel()
			wg.Wait()
			ctx, cancel = context.WithCancel(context.Background())
			di.Clock.Sync()
		}
	}
}
