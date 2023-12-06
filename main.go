package main

import (
	"os"

	"github.com/data-preservation-programs/spade-tenant/api/v1"
	"github.com/data-preservation-programs/spade-tenant/db"
	"github.com/data-preservation-programs/spade-tenant/initializers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func init() {
	initializers.LoadEnvVariables()
	db.ConnectToDB()
}

func main() {
	e := echo.New()
	e.Use(middleware.RequestID())
	api.NewApiV1().RegisterRoutes(e)
	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
