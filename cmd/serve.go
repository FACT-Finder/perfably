package cmd

import (
	"github.com/FACT-Finder/perfably/auth"
	"github.com/FACT-Finder/perfably/config"
	"github.com/FACT-Finder/perfably/router"
	"github.com/FACT-Finder/perfably/server"
	"github.com/FACT-Finder/perfably/state"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func Serve() *cli.Command {
	return &cli.Command{
		Name:  "serve",
		Usage: "start the web service",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "config",
				Aliases:  []string{"c"},
				EnvVars:  []string{"PERFABLY_CONFIG"},
				Usage:    "load configuration from `FILE`",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "address",
				Usage:   "the address to listen on",
				EnvVars: []string{"PERFABLY_ADDRESS"},
				Value:   ":8000",
			},
			&cli.StringFlag{
				Name:    "state",
				Usage:   "the state directory",
				EnvVars: []string{"PERFABLY_STATE"},
				Value:   "./data",
			},
		},
		Action: func(c *cli.Context) error {
			cfg, err := config.New(c.String("config"))
			if err != nil {
				return err
			}
			appState, err := state.ReadState(cfg, c.String("state"))
			if err != nil {
				return err
			}

			users, err := auth.Parse(c.String("state"))
			if err != nil {
				return err
			}
			if err := users.HotReload(); err != nil {
				return err
			}

			handler := router.New(cfg, appState, users)

			listenAddr := c.String("address")
			log.Info().Str("address", listenAddr).Msg("HTTP")
			err = server.Start(handler, listenAddr)

			appState.Close()
			return err
		},
	}
}
