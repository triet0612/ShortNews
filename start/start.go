package start

import (
	"context"
	"log"
	"net/http"
	"newscrapper/internal/config"
	"newscrapper/internal/handler"
	"newscrapper/internal/service"
	"sync"
)

func RunServices(di *config.DI) {
	go runHTTPServer(di)
	runNewsService(di)
}

func runHTTPServer(di *config.DI) {
	h := handler.NewHandler(di.DbCon, di.Clock, di.Signal)
	mux := http.NewServeMux()
	h.Mount(mux)
	server := http.Server{
		Addr:    "localhost:8000",
		Handler: mux,
	}
	log.Fatal(server.ListenAndServe())
}

func runNewsService(di *config.DI) {
	fetcher := service.NewRSSFetchService(di.DbCon)
	audio := service.NewAudioService(di.DbCon, di.LangAudio)
	summary := service.NewSummarizeService(di.DbCon, di.Llmodel, audio)

	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	for {
		select {
		case <-di.Clock.Timer.C:
			di.Clock.Timeout()
			wg.Add(1)
			go func() {
				defer wg.Done()
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
