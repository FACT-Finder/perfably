package main

import (
	"fmt"
	"os"

	"github.com/FACT-Finder/perfably/config"
	"github.com/FACT-Finder/perfably/logger"
	"github.com/FACT-Finder/perfably/router"
	"github.com/FACT-Finder/perfably/server"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	logger.Init(zerolog.ErrorLevel)
	if len(os.Args) != 2 {
		log.Fatal().Msg("Requires on config parameter")
		return
	}

	cfg, err := config.New(os.Args[1])
	if err != nil {
		log.Fatal().Err(err).Msg("Could not read config")
		return
	}

	addr := os.Getenv("PERFABLY_ADDRESS")
	if addr == "" {
		addr = ":8000"
	}

	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("PERFABLY_REDIS_ADDRESS"),
		Password: os.Getenv("PERFABLY_REDIS_PASSWORD"),
	})
	r := router.New(cfg, client)
	fmt.Println("Listening on", addr)
	server.Start(r, addr)
}
