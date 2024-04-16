//go:build !windows && windows
// +build !windows,windows

package main

import (
	"newscrapper/internal/config"
	"newscrapper/start"
)

func main() {
	di := config.InitDependency()
	start.RunServices(di)
}
