package route

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/FACT-Finder/perfably/model"
	"github.com/FACT-Finder/perfably/state"
	"github.com/coreos/go-semver/semver"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func Value(s *state.State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		keys := r.URL.Query()["key"]
		startStr := r.URL.Query().Get("start")
		endStr := r.URL.Query().Get("end")

		project, ok := s.Projects[vars["project"]]
		if !ok {
			badRequest(w, fmt.Sprintf("project not found: %s", vars["project"]))
			return
		}

		project.Lock.RLock()
		defer project.Lock.RUnlock()

		startIndex, endIndex, err := filter(project.Versions, startStr, endStr)
		if err != nil {
			badRequest(w, fmt.Sprintf("could not filter: %s", err))
		}

		result := []model.ReportEntry{}
		for _, v := range project.Versions[startIndex : endIndex+1] {
			vState, ok := project.Data[*v]
			if !ok {
				continue
			}

			entry := model.ReportEntry{
				Key:    *v,
				Values: state.DataPoint{},
				Meta:   vState.Meta,
			}
			for _, key := range keys {
				if value, ok := vState.Values[key]; ok {
					entry.Values[key] = value
				}
			}
			result = append(result, entry)
		}

		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			log.Warn().Err(err).Msg("could not encode to json")
		}
	}
}

func filter(versions semver.Versions, startStr, endStr string) (int, int, error) {
	start, err := semver.NewVersion(startStr)
	if err != nil && startStr != "" {
		return -1, -1, fmt.Errorf("invalid start: %s", err)
	}
	end, err := semver.NewVersion(endStr)
	if err != nil && endStr != "" {
		return -1, -1, fmt.Errorf("invalid end: %s", err)
	}

	startIndex := 0
	if start != nil {
		startIndex = sort.Search(len(versions), func(i int) bool {
			return !versions[i].LessThan(*start)
		})
	}
	endIndex := len(versions) - 1
	if end != nil {
		endIndex = sort.Search(len(versions), func(i int) bool {
			return !versions[i].LessThan(*end)
		})
		if endIndex >= len(versions) {
			endIndex = len(versions) - 1
		}
	}
	return startIndex, endIndex, nil
}
