package main

import (
	"newscrapper/internal/config"
	"newscrapper/service"
)

func main() {

	di := config.InitDependency()

	s := service.RssPullService{Di: di}
	go s.RunRssService()

	service.RunHTTPServer(di)
}
