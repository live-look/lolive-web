package main

import (
	"os"

	"github.com/live-look/lolive-web/internal/server"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:   "lolive-web",
		Action: startServer,
	}

	app.Run(os.Args)
}

func startServer(c *cli.Context) error {
	server := server.AppNew()

	return server.Start()
}
