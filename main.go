package main

import (
	"github.com/data-preservation-programs/spade-tenant/api/v1"
	"github.com/data-preservation-programs/spade-tenant/config"
	"github.com/data-preservation-programs/spade-tenant/db"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	config := config.InitConfig()
	db, err := db.OpenDatabase(config.DB_URL, config.DEBUG, config.DRY_RUN)

	if err != nil {
		e.Logger.Error(err)
	}

	a := api.NewApiV1(db, &config)
	a.RegisterRoutes(e)
	e.Logger.Fatal(e.Start(":" + config.PORT))
}
