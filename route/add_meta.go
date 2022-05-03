package route

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/FACT-Finder/perfably/state"
	"github.com/coreos/go-semver/semver"
	"github.com/gorilla/mux"
)

func AddMeta(s *state.State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		id, err := semver.NewVersion(vars["id"])
		if err != nil {
			badRequest(w, fmt.Sprintf("invalid report id %s: %s", vars["id"], err))
			return
		}

		project, ok := s.Projects[vars["project"]]
		if !ok {
			badRequest(w, fmt.Sprintf("project not found: %s", vars["project"]))
			return
		}

		point, err := parseMetaPoint(r)
		if err != nil {
			badRequest(w, fmt.Sprintf("could not parse request: %s", err))
			return
		}

		project.Lock.Lock()
		defer project.Lock.Unlock()
		project.Add(&state.VersionDataLine{
			Version: *id,
			VersionData: state.VersionData{
				Meta: point,
			},
		})

		w.WriteHeader(http.StatusOK)
		writeString(w, "ok")
	}
}

func parseMetaPoint(r *http.Request) (state.MetaPoint, error) {
	point := state.MetaPoint{}

	if err := json.NewDecoder(r.Body).Decode(&point); err != nil {
		return point, fmt.Errorf("invalid json: %s", err)
	}
	return point, nil
}
