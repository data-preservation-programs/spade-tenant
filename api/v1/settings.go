package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Settings struct {
	SpAutoApprove  bool `json:"sp_auto_approve"`
	SpAutoSuspend  bool `json:"sp_auto_suspend"`
	MaxInFlightGiB uint `json:"max_in_flight_gib"`
}

// handleGetSettings godoc
//
//		@Summary		Apply new Tenant Settings
//		@Security apiKey
//	  @Param 			settings body Settings true "New settings to apply"
//		@Produce		json
//		@Success		200	{object}	ResponseEnvelope{response=Settings}
//		@Router			/settings [post]
func (s *apiV1) handleGetSettings(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{})
}

// handleSetSettings godoc
//
//	@Summary		Get the currently active Tenant Settings
//	@Security apiKey
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=Settings}
//	@Router			/settings [get]
func (s *apiV1) handleSetSettings(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{})
}
