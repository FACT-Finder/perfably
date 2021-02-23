package route

import (
	"io"
	"net/http"
)

func badRequest(w http.ResponseWriter, err string) {
	w.WriteHeader(400)
	io.WriteString(w, err)
}
