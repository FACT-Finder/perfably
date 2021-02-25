package route

import (
	"encoding/json"
	"net/http"

	"github.com/FACT-Finder/perfably/config"
)

func Config(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(cfg)
	}
}
