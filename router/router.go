package router

import (
	"net/http"

	"github.com/FACT-Finder/perfably/config"
	"github.com/FACT-Finder/perfably/route"
	"github.com/FACT-Finder/perfably/state"
	"github.com/FACT-Finder/perfably/auth"
	"github.com/FACT-Finder/perfably/ui"
	"github.com/gorilla/mux"
)

func New(cfg *config.Config, s *state.State, a *auth.Auth) *mux.Router {
	router := mux.NewRouter()
	router.Methods("GET").Path("/project/{project}/value").HandlerFunc(route.Value(s))
	router.Methods("GET").Path("/project/{project}/id").HandlerFunc(route.Ids(s))
	router.Methods("GET").Path("/config").HandlerFunc(route.Config(cfg))
	router.Methods("POST").Path("/project/{project}/report/{id}").HandlerFunc(a.Secure(route.AddReport(s)))
	router.Methods("DELETE").Path("/project/{project}/report/{id}").HandlerFunc(a.Secure(route.DeleteReport(s)))
	router.PathPrefix("/").Handler(AddPrefix("/build", http.FileServer(http.FS(ui.FS))))
	return router
}

func AddPrefix(prefix string, handler http.Handler) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		r.URL.Path = prefix + r.URL.Path
		handler.ServeHTTP(rw, r)
	}
}
