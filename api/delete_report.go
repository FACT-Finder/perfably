package api

import (
	"fmt"
	"net/http"
	"sort"

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

	idx := sort.Search(len(project.Versions), func(i int) bool {
		return !project.Versions[i].LessThan(*id)
	})
	if idx < len(project.Versions) && *project.Versions[idx] == *id {
		delete(project.Data, *id)
		project.Versions = append(project.Versions[:idx], project.Versions[idx+1:]...)
	}
	return ctx.NoContent(http.StatusNoContent)
}
