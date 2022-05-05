package route

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/FACT-Finder/perfably/state"
	"github.com/gorilla/mux"
)

func Ids(state *state.State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		stateProject, ok := state.Projects[vars["project"]]
		if !ok {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("project not found: %s", vars["project"]))
			return
		}

		stateProject.Lock.RLock()
		defer stateProject.Lock.RUnlock()

		err := json.NewEncoder(w).Encode(&stateProject.Versions)
		if err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("could not encode to json: %s", err))
		}
	}
}
