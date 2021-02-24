package route

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"sort"
	"strconv"

	"github.com/FACT-Finder/perfably/config"
	"github.com/FACT-Finder/perfably/model"
	"github.com/FACT-Finder/perfably/rediskey"
	"github.com/coreos/go-semver/semver"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func Value(cfg *config.Config, client *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		keys := r.URL.Query()["key"]
		limitStr := r.URL.Query().Get("limit")
		sortDirection := r.URL.Query().Get("sort")
		startStr := r.URL.Query().Get("start")
		endStr := r.URL.Query().Get("end")

		project, ok := cfg.Projects[vars["project"]]
		if !ok {
			badRequest(w, fmt.Sprintf("project not found: %s", vars["project"]))
			return
		}

		reportsKey := rediskey.ReportIDs(vars["project"])
		ids, err := client.SMembers(context.Background(), reportsKey).Result()
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			writeString(w, fmt.Sprintf("redis get failed: %s", err))
			return
		}

		if len(ids) == 0 {
			_ = json.NewEncoder(w).Encode([]model.ReportEntry{})
			return
		}

		var sortedFilteredReportIds []string
		if project.IDType == config.ReportIDTypeSemver {
			sortedFilteredReportIds, err = sortedFilteredSemverSlice(ids, startStr, endStr, sortDirection)
		} else {
			sortedFilteredReportIds, err = sortedFilteredIntSlice(ids, startStr, endStr, sortDirection)
		}

		if (err != nil) {
			badRequest(w, fmt.Sprintf("could not parse ids: %s", err))
			return
		}

		limit, _ := strconv.Atoi(limitStr)
		if limit <= 0 {
			limit = len(ids)
		}

		reportIds := sortedFilteredReportIds[:int(math.Min(float64(len(sortedFilteredReportIds)), float64(limit)))]

		metrics := []string{}

		for _, reportID := range reportIds {
			for _, key := range keys {
				metrics = append(metrics, rediskey.Metric(vars["project"], reportID, key))
			}
		}

		values, err := client.MGet(context.Background(), metrics...).Result()
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			writeString(w, fmt.Sprintf("redis get failed: %s", err))
			return
		}

		result := []model.ReportEntry{}
		for i, reportID := range reportIds {
			entry := model.ReportEntry{Key: reportID, Values: model.DataPoint{}}
			for j, metric := range keys {
				x := i*len(keys) + j
				if values[x] != nil {
					value, _ := strconv.ParseFloat(values[x].(string), 64)
					entry.Values[metric] = value
				}
			}
			result = append(result, entry)
		}

		err = json.NewEncoder(w).Encode(result)
		if (err != nil) {
			log.Error().Err(err).Msg("could not encode to json")
		}
	}
}

func sortedFilteredIntSlice(ids []string, startStr, endStr, sortDirection string) ([]string, error) {
	start, err := strconv.ParseInt(startStr, 10, 64)
	if err != nil && startStr != "" {
		return nil, fmt.Errorf("invalid start: %s", err)
	}
	end, err := strconv.ParseInt(endStr, 10, 64)
	if err != nil && endStr != "" {
		return nil, fmt.Errorf("invalid end: %s", err)
	}

	asInts := []int64{}

	for _, id := range ids {
		i, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return nil, err
		}
		asInts = append(asInts, i)
	}

	if sortDirection == "desc" {
		sort.Slice(asInts, func(i, j int) bool {
			return asInts[i] > asInts[j]
		})
	} else {
		sort.Slice(ids, func(i, j int) bool {
			return asInts[i] < asInts[j]
		})
	}

	result := []string{}

	for _, i := range asInts {
		if (endStr == "" || i <= end) && (startStr == "" || i >= start) {
			result = append(result, fmt.Sprint(i))
		}
	}

	return result, nil
}

func sortedFilteredSemverSlice(ids []string, startStr, endStr, sortDirection string) ([]string, error) {
	start, err := semver.NewVersion(startStr)
	if err != nil && startStr != "" {
		return nil, fmt.Errorf("invalid start: %s", err)
	}
	end, err := semver.NewVersion(endStr)
	if err != nil && endStr != "" {
		return nil, fmt.Errorf("invalid end: %s", err)
	}

	asSemvers := semver.Versions{}

	for _, id := range ids {
		i, err := semver.NewVersion(id)
		if err != nil {
			return nil, err
		}
		asSemvers = append(asSemvers, i)
	}

	if sortDirection == "desc" {
		sort.Sort(sort.Reverse(asSemvers))
	} else {
		sort.Sort(asSemvers)
	}

	result := []string{}

	for _, version := range asSemvers {
		if (end == nil || !end.LessThan(*version)) && (start == nil || !version.LessThan(*start)) {
			result = append(result, version.String())
		}
	}

	return result, nil
}
