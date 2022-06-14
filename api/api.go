package api

import (
	"fmt"
	"net/http"

	"github.com/FACT-Finder/perfably/auth"
	"github.com/FACT-Finder/perfably/config"
	"github.com/FACT-Finder/perfably/state"
	"github.com/FACT-Finder/perfably/swagger"
	"github.com/FACT-Finder/perfably/ui"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func New(cfg *config.Config, s *state.State, users *auth.Auth) *echo.Echo {
	app := echo.New()
	app.Use(middleware.Recover())
	app.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		message := ""
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			message = fmt.Sprint(he.Message)
		} else {
			message = err.Error()
		}
		c.JSON(code, &ApiError{
			Error:       http.StatusText(code),
			Description: message,
		})
	}

	wrapper := ServerInterfaceWrapper{Handler: &api{config: cfg, s: s}}

	app.GET("/config", wrapper.GetConfig)
	app.GET("/project/:project/id", wrapper.GetIds)
	app.DELETE("/project/:project/report/:version", users.Secure(wrapper.DeleteReport))
	app.POST("/project/:project/report/:version", users.Secure(wrapper.AddMetrics))
	app.POST("/project/:project/report/:version/meta", users.Secure(wrapper.AddMeta))
	app.GET("/project/:project/value", wrapper.GetValues)
	app.StaticFS("/", ui.Build)

	spec, err := GetSwagger()
	if err != nil {
		panic(err)
	}
	swagger.Register(app, spec)

	return app
}

type api struct {
	config *config.Config
	s      *state.State
}
