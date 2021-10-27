package main

import (
	"github.com/live-look/lolive-web/internal/server"
)

func main() {
	app := &server.App{}

	app.Start()
}
