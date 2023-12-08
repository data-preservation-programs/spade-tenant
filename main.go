package main

import (
	"github.com/data-preservation-programs/spade-tenant/api/v1"
	"github.com/data-preservation-programs/spade-tenant/config"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	config := config.InitConfig()
	api.RegisterRoutes(e, config)
	e.Logger.Fatal(e.Start(":" + config.PORT))
}
