package main

import (
	"log"
	"os"

	"github.com/hulutech-web/goravel-kit-cli/internal/commands"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:     "goravel-kit-cli",
		Usage:    "A CLI tool to create new Goravel applications from templates",
		Version:  "v1.0.0",
		Commands: []*cli.Command{commands.NewCommand},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
