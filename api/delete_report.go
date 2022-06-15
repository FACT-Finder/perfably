package api

import (
	"fmt"
	"net/http"

	"github.com/FACT-Finder/perfably/state"
	"github.com/coreos/go-semver/semver"
	"github.com/labstack/echo/v4"
)

func (a *api) DeleteReport(ctx echo.Context, projectName, version string) error {
	project, ok := a.s.Projects[projectName]
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("project not found: %s", projectName))
	}

	id, err := semver.NewVersion(version)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid report id %s: %s", version, err))
	}

	project.Lock.Lock()
	defer project.Lock.Unlock()

	project.Add(&state.VersionDataLine{
		Version: *id,
		Delete:  true,
	})
	return ctx.NoContent(http.StatusNoContent)
}
