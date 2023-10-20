package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type TenantSettings struct {
	AutoApprove bool `json:"auto_approve"`
	AutoSuspend bool `json:"auto_suspend"`
}

// handleSetTenantSettings godoc
//	@Summary		Apply new Tenant Settings
// 	@Param 		  token header string true "Auth token"
//  @Param 			settings body TenantSettings true "New settings to apply"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{Response=TenantSettings}
//	@Router			/tenant/settings [post]
func (s *apiV1) handleSetTenantSettings(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{})
}

// handleGetTenantSettings godoc
//	@Summary		Get the currently active Tenant Settings
// 	@Param 		  token header string true "Auth token"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{Response=TenantSettings}
//	@Router			/tenant/settings [get]
func (s *apiV1) handleGetTenantSettings(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{})
}
