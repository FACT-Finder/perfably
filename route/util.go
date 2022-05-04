package route

import (
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

func badRequest(w http.ResponseWriter, err string) {
	w.WriteHeader(400)
	writeString(w, err)
}

func internalServerError(w http.ResponseWriter, err string) {
	w.WriteHeader(500)
	writeString(w, err)
}

func writeString(w io.Writer, msg string) {
	_, err := io.WriteString(w, msg)
	if err != nil {
		log.Warn().Err(err).Msg("could not write response")
	}
}
