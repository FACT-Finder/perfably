package main

import (
	"fmt"
	"os"

	"github.com/FACT-Finder/perfably/cmd"
	"github.com/FACT-Finder/perfably/logger"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

func main() {
	logger.Init(zerolog.ErrorLevel)
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "redis-address",
				Usage:   "address of redis server",
				EnvVars: []string{"PERFABLY_REDIS_ADDRESS"},
				Value:   ":6379",
			},
			&cli.StringFlag{
				Name:    "redis-password",
				EnvVars: []string{"PERFABLY_REDIS_PASSWORD"},
				Usage:   "password of redis server",
			},
		},
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
