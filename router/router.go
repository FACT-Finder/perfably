package router

import (
	"github.com/FACT-Finder/perfably/config"
	"github.com/FACT-Finder/perfably/route"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

func New(cfg *config.Config, client *redis.Client) *mux.Router {
	router := mux.NewRouter()
	router.Methods("GET").Path("/project/{project}/value").HandlerFunc(route.Value(cfg, client))
	router.Methods("GET").Path("/project/{project}/id").HandlerFunc(route.Ids(cfg, client))
	router.Methods("GET").Path("/project/{project}/metrics").HandlerFunc(route.Metrics(cfg, client))
	router.Methods("POST").Path("/project/{project}/report/{id}").HandlerFunc(route.AddReport(cfg, client))
	return router
}
