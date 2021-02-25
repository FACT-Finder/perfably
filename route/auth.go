package route

import (
	"context"
	"fmt"
	"net/http"

	"github.com/FACT-Finder/perfably/rediskey"
	"github.com/FACT-Finder/perfably/token"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

func BasicAuth(handler http.HandlerFunc, client *redis.Client, realm string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()

		if !ok {
			unauthorized(w, realm)
			return
		}

		hash, err := client.HGet(context.Background(), rediskey.Tokens(), user).Result()
		if err != nil {
			log.Error().Err(err).Msg("redis")
			w.WriteHeader(http.StatusBadGateway)
			writeString(w, fmt.Sprintf("redis failed: %s", err))
			return
		}

		if !token.ComparePassword([]byte(hash), []byte(pass)) {
			unauthorized(w, realm)
			return
		}

		handler(w, r)
	}
}

func unauthorized(w http.ResponseWriter, realm string) {
	w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
	w.WriteHeader(401)
	writeString(w, "unauthorized")
}
