package route

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/FACT-Finder/perfably/state"
	"github.com/coreos/go-semver/semver"
	"github.com/gorilla/mux"
)

func AddReport(s *state.State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		id, err := semver.NewVersion(vars["id"])
		if err != nil {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid report id %s: %s", vars["id"], err))
			return
		}

		project, ok := s.Projects[vars["project"]]
		if !ok {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("project not found: %s", vars["project"]))
			return
		}

		point, err := parseDataPoint(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("could not parse request: %s", err))
			return
		}

		project.Lock.Lock()
		defer project.Lock.Unlock()

		err = project.Add(&state.VersionDataLine{
			Version: *id,
			VersionData: state.VersionData{
				Values: point,
			},
		})
		if err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("could not add line: %s", err))
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func parseDataPoint(r *http.Request) (state.DataPoint, error) {
	point := state.DataPoint{}

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
