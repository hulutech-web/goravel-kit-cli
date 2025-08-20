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
		Description: `Goravel Kit CLI - Quickly create new Goravel projects from template.

Examples:
  goravel-kit-cli new my-app
  goravel-kit-cli new my-app --ssh --verbose
  goravel-kit-cli new my-app --branch develop --force`,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
