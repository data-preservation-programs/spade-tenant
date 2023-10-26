package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Settings struct {
	AutoApprove bool `json:"auto_approve"`
	AutoSuspend bool `json:"auto_suspend"`
}

// handleGetSettings godoc
//	@Summary		Apply new Tenant Settings
// 	@Param 		  token header string true "Auth token"
//  @Param 			settings body Settings true "New settings to apply"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=Settings}
//	@Router			/settings [post]
func (s *apiV1) handleGetSettings(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{})
}

// handleSetSettings godoc
//	@Summary		Get the currently active Tenant Settings
// 	@Param 		  token header string true "Auth token"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=Settings}
//	@Router			/settings [get]
func (s *apiV1) handleSetSettings(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{})
}
