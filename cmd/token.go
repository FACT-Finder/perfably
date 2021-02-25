package cmd

import (
	"context"
	"fmt"

	"github.com/FACT-Finder/perfably/rediskey"
	"github.com/FACT-Finder/perfably/token"
	"github.com/go-redis/redis/v8"
	"github.com/urfave/cli/v2"
)

const (
	tokenLength      = 32
	passwordStrength = 12
)

func Token() *cli.Command {
	return &cli.Command{
		Name:  "token",
		Usage: "manage tokens",
		Subcommands: []*cli.Command{
			{
				Name:  "list",
				Usage: "list tokens",
				Action: func(c *cli.Context) error {
					addr := c.String("redis-address")
					client := redis.NewClient(&redis.Options{
						Addr:     addr,
						Password: c.String("redis-password"),
					})
					return listTokens(client)
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
					addr := c.String("redis-address")
					client := redis.NewClient(&redis.Options{
						Addr:     addr,
						Password: c.String("redis-password"),
					})

					return createToken(client, c.String("name"))
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
					addr := c.String("redis-address")
					client := redis.NewClient(&redis.Options{
						Addr:     addr,
						Password: c.String("redis-password"),
					})

					return removeToken(client, c.String("name"))
				},
			},
		},
	}
}

func listTokens(client *redis.Client) error {
	names, err := client.HKeys(context.Background(), rediskey.Tokens()).Result()
	if err != nil {
		return fmt.Errorf("could not read from redis: %s", err)
	}

	for _, name := range names {
		fmt.Println(name)
	}
	return nil
}

func createToken(client *redis.Client, name string) error {
	password := token.GenerateRandomString(tokenLength)
	hashedPassword := token.CreatePassword(password, passwordStrength)

	created, err := client.HSetNX(context.Background(), rediskey.Tokens(), name, hashedPassword).Result()
	if err != nil {
		return fmt.Errorf("could not write to redis: %s", err)
	}

	if !created {
		return fmt.Errorf("token '%s' already exists", name)
	}

	fmt.Printf("%s:%s\n", name, password)
	return nil
}

func removeToken(client *redis.Client, name string) error {
	count, err := client.HDel(context.Background(), rediskey.Tokens(), name).Result()
	if err != nil {
		return fmt.Errorf("could not delete token: %s", err)
	}

	fmt.Printf("%d token removed\n", count)
	return nil
}
