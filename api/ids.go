package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (a *api) GetIds(ctx echo.Context, projectName string) error {
	stateProject, ok := a.s.Projects[projectName]
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("project not found: %s", projectName))
	}

	stateProject.Lock.RLock()
	defer stateProject.Lock.RUnlock()

	return ctx.JSON(http.StatusOK, &stateProject.Versions)
}
