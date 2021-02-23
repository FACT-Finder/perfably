package main

import (
	"fmt"
	"log"
	"os"

	"github.com/FACT-Finder/perfably/config"
	"github.com/FACT-Finder/perfably/router"
	"github.com/FACT-Finder/perfably/server"
	"github.com/go-redis/redis/v8"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Requires on config parameter")
		return
	}

	cfg, err := config.New(os.Args[1])
	if err != nil {
		log.Fatalln("Could not read config", err)
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
