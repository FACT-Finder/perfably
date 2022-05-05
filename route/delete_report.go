package route

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/FACT-Finder/perfably/state"
	"github.com/coreos/go-semver/semver"
	"github.com/gorilla/mux"
)

func DeleteReport(s *state.State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		project, ok := s.Projects[vars["project"]]
		if !ok {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("project not found: %s", vars["project"]))
			return
		}

		id, err := semver.NewVersion(vars["id"])
		if err != nil {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid report id %s: %s", vars["id"], err))
			return
		}

		project.Lock.Lock()
		defer project.Lock.Unlock()

		idx := sort.Search(len(project.Versions), func(i int) bool {
			return !project.Versions[i].LessThan(*id)
		})
		if idx < len(project.Versions) && *project.Versions[idx] == *id {
			delete(project.Data, *id)
			project.Versions = append(project.Versions[:idx], project.Versions[idx+1:]...)
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
