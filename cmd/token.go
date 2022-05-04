package cmd

import (
	"fmt"

	"github.com/FACT-Finder/perfably/auth"
	"github.com/urfave/cli/v2"
)

func Token() *cli.Command {
	return &cli.Command{
		Name:  "token",
		Usage: "manage tokens",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "state",
				Usage:   "the state directory",
				EnvVars: []string{"PERFABLY_STATE"},
				Value:   "data",
			},
		},
		Subcommands: []*cli.Command{
			{
				Name:  "list",
				Usage: "list tokens",
				Action: func(c *cli.Context) error {
					a, err := auth.Parse(c.String("state"))
					if err != nil {
						return err
					}
					for _, name := range a.Names() {
						fmt.Println(name)
					}
					return nil
				},
			},
			{
				Name:  "create",
				Usage: "create a new token",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "name",
						Aliases:  []string{"n"},
						Usage:    "token name",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					name := c.String("name")
					a, err := auth.Parse(c.String("state"))
					if err != nil {
						return err
					}
					hash, err := a.Create(c.String("name"))
					if err != nil {
						return err
					}
					fmt.Printf("%s:%s\n", name, hash)
					return nil
				},
			},
			{
				Name:  "rm",
				Usage: "remove a token",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "name",
						Aliases:  []string{"n"},
						Usage:    "token name",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					a, err := auth.Parse(c.String("state"))
					if err != nil {
						return err
					}
					return a.Remove(c.String("name"))
				},
			},
		},
	}
}
