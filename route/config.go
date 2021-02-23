package route

import (
	"encoding/json"
	"github.com/FACT-Finder/perfably/config"
	"net/http"
)

func Config(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(cfg)
	}
}
