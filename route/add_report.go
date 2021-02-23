package route

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/FACT-Finder/perfably/config"
	"github.com/FACT-Finder/perfably/model"
	"github.com/FACT-Finder/perfably/rediskey"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

func AddReport(cfg *config.Config, client *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		project, ok := cfg.Projects[vars["project"]]
		if !ok {
			badRequest(w, fmt.Sprintf("project not found: %s", vars["project"]))
			return
		}

		if err := model.ValidId(project.IDType, vars["id"]); err != nil {
			badRequest(w, fmt.Sprintf("invalid report id %s must be a %s", vars["id"], project.IDType))
			return
		}
		point, err := parseDataPoint(r)
		if err != nil {
			badRequest(w, fmt.Sprintf("could not parse request: %s", err))
			return
		}

		metrics := []interface{}{}
		for key := range point {
			metrics = append(metrics, key)
			layers := strings.Split(key, ".")
			if len(layers) != len(project.Layers) {
				badRequest(w, fmt.Sprintf("%s must have %d layers", key, len(project.Layers)))
				return
			}
		}

		pipe := client.TxPipeline()

		for metricID, value := range point {
			metricKey := rediskey.Metric(vars["project"], vars["id"], metricID)
			pipe.Set(context.Background(), metricKey, value, 0)
		}

		reportsKey := rediskey.ReportIDs(vars["project"])
		pipe.SAdd(context.Background(), reportsKey, vars["id"])

		metricsKey := rediskey.Metrics(vars["project"])
		pipe.SAdd(context.Background(), metricsKey, metrics...)

		if _, err := pipe.Exec(context.Background()); err != nil {
			w.WriteHeader(http.StatusBadGateway)
			writeString(w, fmt.Sprintf("redis failed: %s", err))
			return
		}

		w.WriteHeader(http.StatusOK)
		writeString(w, "ok")
	}
}

func parseDataPoint(r *http.Request) (model.DataPoint, error) {
	point := model.DataPoint{}

	if strings.Contains(r.Header.Get("content-type"), "application/x-www-form-urlencoded") {
		for key, values := range r.PostForm {
			if len(values) != 1 {
				return point, fmt.Errorf("key %s has multiple values", key)
			}
			var err error
			point[key], err = strconv.ParseFloat(values[0], 64)
			if err != nil {
				return point, fmt.Errorf("invalid value %s in key %s", values[0], key)
			}
		}
	} else {
		if err := json.NewDecoder(r.Body).Decode(&point); err != nil {
			return point, fmt.Errorf("invalid json: %s", err)
		}
	}
	return point, nil
}
