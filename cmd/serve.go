package cmd

import (
	"fmt"

	"github.com/FACT-Finder/perfably/config"
	"github.com/FACT-Finder/perfably/router"
	"github.com/FACT-Finder/perfably/server"
	"github.com/go-redis/redis/v8"
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
				Usage:   "perfably listen address",
				EnvVars: []string{"PERFABLY_ADDRESS"},
				Value:   ":8000",
			},
		},
		Action: func(c *cli.Context) error {
			cfg, err := config.New(c.String("config"))
			if err != nil {
				return err
			}

			redisAddr := c.String("redis-address")
			client := redis.NewClient(&redis.Options{
				Addr:     redisAddr,
				Password: c.String("redis-password"),
			})
			r := router.New(cfg, client)

			listenAddr := c.String("address")
			fmt.Println("Listening on", listenAddr)
			return server.Start(r, listenAddr)
		},
	}
}
