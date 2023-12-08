package api

import (
	"net/http"

	"github.com/data-preservation-programs/spade-tenant/config"
	"github.com/labstack/echo/v4"
)

// handleStatus godoc
//
//	@Summary		Simple health check endpoint
//	@Description	This endpoint is used to check the health of the service
//	@Security apiKey
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=string}
//	@Router			/status [get]
func ConfigureStatusRouter(e *echo.Group, config config.TenantServiceConfig) {
	status := e.Group("/status")

	status.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, CreateSuccessResponseEnvelope(c, "Healthy"))
	})
}
