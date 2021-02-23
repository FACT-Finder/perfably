package route

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/FACT-Finder/perfably/config"
	"github.com/FACT-Finder/perfably/rediskey"
	"github.com/coreos/go-semver/semver"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

func Ids(cfg *config.Config, client *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		limitStr := r.URL.Query().Get("limit")
		query := r.URL.Query().Get("query")

		project, ok := cfg.Projects[vars["project"]]
		if !ok {
			badRequest(w, fmt.Sprintf("project not found: %s", vars["project"]))
			return
		}

		reportsKey := rediskey.ReportIDs(vars["project"])
		ids, err := client.SMembers(context.Background(), reportsKey).Result()
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			io.WriteString(w, fmt.Sprintf("redis get failed: %s", err))
			return
		}
		filteredIds := []string{}
		for _, id := range ids {
			if strings.HasPrefix(id, query) {
				filteredIds = append(filteredIds, id)
			}
		}

		var sortedIds []string
		if project.IDType == config.ReportIDTypeInt {
			sortedIds, err = sortedInt(filteredIds)
		} else {
			sortedIds, err = sortedSemver(filteredIds)
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, fmt.Sprintf("could not convert ids: %s", err))
			return
		}

		limit, _ := strconv.Atoi(limitStr)
		if limit <= 0 {
			limit = len(sortedIds)
		}

		json.NewEncoder(w).Encode(sortedIds[:int(math.Min(float64(len(sortedIds)), float64(limit)))])
	}
}

func toSemvers(ids []string) (semver.Versions, error) {
	versions := semver.Versions{}
	for _, id := range ids {
		if ver, err := toSemver(id); err == nil {
			versions = append(versions, ver)
		} else {
			return versions, fmt.Errorf("could not parse: %s", id)
		}
	}
	return versions, nil
}

func toSemver(id string) (*semver.Version, error) {
	if ver, err := semver.NewVersion(id); err == nil {
		return ver, nil
	}
	if _, err := strconv.ParseInt(id, 10, 64); err == nil {
		return &semver.Version{PreRelease: semver.PreRelease(id)}, nil
	}
	return nil, fmt.Errorf("could not parse: %s", id)
}

func toStringSemver(version *semver.Version) string {
	if version.Major == 0 && version.Minor == 0 && version.Patch == 0 {
		return string(version.PreRelease)
	}
	return version.String()
}

func sortedInt(ids []string) ([]string, error) {
	asInts := []int64{}

	for _, id := range ids {
		i, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return nil, err
		}
		asInts = append(asInts, i)
	}

	sort.Slice(ids, func(i, j int) bool {
		return asInts[i] < asInts[j]
	})

	result := []string{}

	for _, i := range asInts {
		result = append(result, fmt.Sprint(i))
	}

	return result, nil
}

func sortedSemver(ids []string) ([]string, error) {
	asSemvers := semver.Versions{}

	for _, id := range ids {
		i, err := semver.NewVersion(id)
		if err != nil {
			return nil, err
		}
		asSemvers = append(asSemvers, i)
	}

	sort.Sort(asSemvers)

	result := []string{}

	for _, version := range asSemvers {
		result = append(result, version.String())
	}

	return result, nil
}
