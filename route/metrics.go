package route

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/FACT-Finder/perfably/config"
	"github.com/FACT-Finder/perfably/rediskey"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

func Metrics(cfg *config.Config, client *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		project, ok := cfg.Projects[vars["project"]]
		if !ok {
			badRequest(w, fmt.Sprintf("project not found: %s", vars["project"]))
			return
		}

		metricsKey := rediskey.Metrics(vars["project"])
		metrics, err := client.SMembers(context.Background(), metricsKey).Result()
		if err != nil {
			badRequest(w, fmt.Sprintf("redis get failed: %s", err))
			return
		}

		result := map[string]interface{}{}
		for _, metric := range metrics {
			layers := strings.Split(metric, ".")
			if len(layers) != len(project.Layers) {
				w.WriteHeader(http.StatusInternalServerError)
				writeString(w, fmt.Sprintf("redis key %s have %d layers", metric, len(project.Layers)))
				return
			}

			part := result
			for i, group := range layers {
				if i == len(layers)-1 {
					part[group] = true
					continue
				}
				subPart, _ := part[group].(map[string]interface{})
				if subPart == nil {
					subPart = map[string]interface{}{}
				}
				part[group] = subPart
				part = subPart
			}
		}
		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			writeString(w, fmt.Sprintf("could not encode to json: %s", err))
		}
	}
}
