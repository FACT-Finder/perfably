package route

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

type Error struct {
	Error       string `json:"error"`
	Description string `json:"description"`
}

func writeError(w http.ResponseWriter, status int, description string) {
	w.WriteHeader(status)
	msg := Error{
		Error:       http.StatusText(status),
		Description: description,
	}
	if err := json.NewEncoder(w).Encode(&msg); err != nil {
		log.Debug().Err(err).Msg("could not write response")
	}
}
