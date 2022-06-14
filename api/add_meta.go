package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/FACT-Finder/perfably/state"
	"github.com/coreos/go-semver/semver"
	"github.com/labstack/echo/v4"
)

func (a *api) AddMeta(ctx echo.Context, projectName, version string) error {
	project, ok := a.s.Projects[projectName]
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("project not found: %s", projectName))
	}

	id, err := semver.NewVersion(version)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid report id %s: %s", version, err))
	}

	point, err := parseMetaPoint(ctx.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("could not parse request: %s", err))
	}
	project.Lock.Lock()
	defer project.Lock.Unlock()
	project.Add(&state.VersionDataLine{
		Version: *id,
		VersionData: state.VersionData{
			Meta: point,
		},
	})

	return ctx.NoContent(http.StatusNoContent)
}

func parseMetaPoint(r *http.Request) (state.MetaPoint, error) {
	point := state.MetaPoint{}

	if err := json.NewDecoder(r.Body).Decode(&point); err != nil {
		return point, fmt.Errorf("invalid json: %s", err)
	}
	return point, nil
}
