package route

import (
	"context"
	"fmt"
	"net/http"

	"github.com/FACT-Finder/perfably/config"
	"github.com/FACT-Finder/perfably/model"
	"github.com/FACT-Finder/perfably/rediskey"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

func DeleteReport(cfg *config.Config, client *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		project, ok := cfg.Projects[vars["project"]]
		if !ok {
			badRequest(w, fmt.Sprintf("project not found: %s", vars["project"]))
			return
		}

		if err := model.ValidID(project.IDType, vars["id"]); err != nil {
			badRequest(w, fmt.Sprintf("invalid report id %s must be a %s", vars["id"], project.IDType))
			return
		}

		metricsKey := rediskey.Metrics(vars["project"])
		keys, err := client.SMembers(context.Background(), metricsKey).Result()
		if err != nil {
			badRequest(w, fmt.Sprintf("redis get failed: %s", err))
			return
		}

		pipe := client.TxPipeline()

		metrics := []string{}
		for _, key := range keys {
			metrics = append(metrics, rediskey.Metric(vars["project"], vars["id"], key))
		}
		pipe.Del(context.Background(), metrics...)

		reportsKey := rediskey.ReportIDs(vars["project"])
		pipe.SRem(context.Background(), reportsKey, vars["id"])

		if _, err := pipe.Exec(context.Background()); err != nil {
			w.WriteHeader(http.StatusBadGateway)
			writeString(w, fmt.Sprintf("redis failed: %s", err))
			return
		}

		w.WriteHeader(http.StatusOK)
		writeString(w, "ok")
	}
}
