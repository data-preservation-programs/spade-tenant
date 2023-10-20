package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Legend []LegendEntry

type LegendEntry struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// handleGetLegend godoc
//	@Summary		Get legend for a tenant
// 	@Param 		  token header string true "Auth token"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{Response=Legend}
//	@Router			/legend [get]
func (s *apiV1) handleGetLegend(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{})
}

// handleSetLegend godoc
//	@Summary		Update legend for a tenant
// 	@Param 		  token header string true "Auth token"
//  @Param 			legend body Legend true "New legend to update"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{Response=Legend}
//	@Router			/legend [post]
func (s *apiV1) handleSetLegend(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{})
}
