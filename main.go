package main

import (
	"github.com/data-preservation-programs/spade-tenant/api/v1"
	"github.com/data-preservation-programs/spade-tenant/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	config := config.InitConfig()
	e.Use(middleware.RequestID())
	api.NewApiV1().RegisterRoutes(e, config)
	e.Logger.Fatal(e.Start(":" + config.Settings.PORT))
}
