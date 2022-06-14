package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/FACT-Finder/perfably/state"
	"github.com/coreos/go-semver/semver"
	"github.com/labstack/echo/v4"
)

func (a *api) AddMetrics(ctx echo.Context, projectName, version string) error {
	project, ok := a.s.Projects[projectName]
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("project not found: %s", projectName))
	}

	id, err := semver.NewVersion(version)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid report id %s: %s", version, err))
	}

	point, err := parseDataPoint(ctx.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("could not parse request: %s", err))
	}

	project.Lock.Lock()
	defer project.Lock.Unlock()

	err = project.Add(&state.VersionDataLine{
		Version:     *id,
		VersionData: state.VersionData{Values: point},
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("could persist metric line: %s", err))
	}

	return ctx.NoContent(http.StatusNoContent)
}

func parseDataPoint(r *http.Request) (state.DataPoint, error) {
	point := state.DataPoint{}

	if strings.Contains(r.Header.Get("content-type"), "application/x-www-form-urlencoded") {
		for key, values := range r.PostForm {
			if len(values) != 1 {
				return point, fmt.Errorf("key %s has multiple values", key)
			}
			var err error
			point[key], err = strconv.ParseFloat(values[0], 64)
			if err != nil {
				return point, fmt.Errorf("invalid value %s in key %s", values[0], key)
			}
		}
	} else {
		if err := json.NewDecoder(r.Body).Decode(&point); err != nil {
			return point, fmt.Errorf("invalid json: %s", err)
		}
	}
	return point, nil
}
