package start

import (
	"context"
	"log"
	"net/http"
	"newscrapper/internal/config"
	"newscrapper/internal/handler"
	"newscrapper/internal/service"
	"os/exec"
	"runtime"
	"sync"
)

func RunServices(di *config.DI) {
	go runNewsService(di)
	runHTTPServer(di)
}

func runHTTPServer(di *config.DI) {
	log.Println("start http server")
	h := handler.NewHandler(di.DbCon, di.Clock, di.Signal)
	mux := http.NewServeMux()
	h.Mount(mux)
	server := http.Server{
		Addr:    ":8000",
		Handler: mux,
	}
	open("http://localhost:8000")
	log.Fatal(server.ListenAndServe())
}

func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func runNewsService(di *config.DI) {
	fetcher := service.NewRSSFetchService(di.DbCon)
	audio := service.NewAudioService(di.DbCon)
	summary := service.NewSummarizeService(di.DbCon, audio)
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
