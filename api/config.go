package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (a *api) GetConfig(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, a.config)
}
