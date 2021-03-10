package router

import (
	"net/http"

	"github.com/FACT-Finder/perfably/config"
	"github.com/FACT-Finder/perfably/route"
	"github.com/FACT-Finder/perfably/ui"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

func New(cfg *config.Config, client *redis.Client) *mux.Router {
	router := mux.NewRouter()
	router.Methods("GET").Path("/project/{project}/value").HandlerFunc(route.Value(cfg, client))
	router.Methods("GET").Path("/project/{project}/id").HandlerFunc(route.Ids(cfg, client))
	router.Methods("GET").Path("/project/{project}/metrics").HandlerFunc(route.Metrics(cfg, client))
	router.Methods("GET").Path("/config").HandlerFunc(route.Config(cfg))
	router.Methods("POST").Path("/project/{project}/report/{id}").HandlerFunc(route.BasicAuth(route.AddReport(cfg, client), client, "perfably"))
	router.Methods("DELETE").Path("/project/{project}/report/{id}").HandlerFunc(route.BasicAuth(route.DeleteReport(cfg, client), client, "perfably"))
	router.PathPrefix("/").Handler(AddPrefix("/build", http.FileServer(http.FS(ui.FS))))
	return router
}

func AddPrefix(prefix string, handler http.Handler) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		r.URL.Path = prefix + r.URL.Path
		handler.ServeHTTP(rw, r)
	}
}
