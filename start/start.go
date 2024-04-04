package start

import (
	"log"
	"net/http"
	"newscrapper/internal/config"
	"newscrapper/internal/handler"
	"newscrapper/internal/service"
)

func RunServices(di *config.DI) {
	go runHTTPServer(di)
	runNewsService(di)
}

func runHTTPServer(di *config.DI) {
	h := handler.NewHandler(di.DbCon, di.Clock)
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
	thumbnail := service.NewThumbnailService(di.DbCon)
	summary := service.NewSummarizeService(di.DbCon, di.Llmodel)
	audio := service.NewAudioService(di.DbCon)

	for {
		<-di.Clock.Timer.C
		di.Clock.Timeout()
		fetcher.NewsExtraction()
		go thumbnail.UpdateThumbnail()
		summary.ArticleSummarize()
		go audio.GenerateAudio()
	}
}
