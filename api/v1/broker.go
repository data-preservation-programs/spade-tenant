package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (a *apiV1) ConfigureBrokerRouter(e *echo.Group) {
	g := e.Group("/broker")
	g.GET("", a.handleGetTenantsInformation)
	g.POST("", a.handlePostNotifyTenantService)
}

// handleGetTenantsInformation godoc
//
//	@Summary		List of all tenants in a JSON object to be consumed by the broker.
//	@Security apiKey
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=LabelsResponse}
//	@Router			/broker [get]
func (a *apiV1) handleGetTenantsInformation(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, CreateErrorResponseEnvelope(c, http.StatusNotImplemented, ""))
}

// handlePostNotifyTenantService godoc
//
//	@Summary		Allows the broker to notify the tenant service.
//	@Security apiKey
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=LabelsResponse}
//	@Router			/broker [post]
func (a *apiV1) handlePostNotifyTenantService(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, CreateErrorResponseEnvelope(c, http.StatusNotImplemented, ""))
}
