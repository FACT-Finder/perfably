package main

import (
	"fmt"
	"os"

	"github.com/FACT-Finder/perfably/cmd"
	"github.com/FACT-Finder/perfably/logger"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

//go:generate oapi-codegen -config ./openapi-gen.yml openapi.yaml
func main() {
	logger.Init(zerolog.InfoLevel)
	app := &cli.App{
		Commands: []*cli.Command{
			cmd.Serve(),
			cmd.Token(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
