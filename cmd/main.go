package main

import (
	"newscrapper/internal/config"
	"newscrapper/start"
	"os/exec"
	"runtime"
)

func main() {
	ensureInstallation()
	di := config.InitDependency()
	open("http://localhost:8000")
	start.RunServices(di)
}

func ensureInstallation() {
	exec.Command("pip", "install", "mycroft-mimic3-tts[all]").Run()
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
	default:
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
