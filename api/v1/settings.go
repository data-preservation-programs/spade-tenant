package api

import (
	"encoding/json"
	"net/http"

	"github.com/data-preservation-programs/spade-tenant/db"
	"github.com/labstack/echo/v4"
)

type Settings struct {
	SpAutoApprove  bool `json:"sp_auto_approve"`
	SpAutoSuspend  bool `json:"sp_auto_suspend"`
	MaxInFlightGiB uint `json:"max_in_flight_gib"`
}

func (a *apiV1) ConfigureSettingsRouter(e *echo.Group) {
	g := e.Group("/settings")
	g.POST("", a.handleSetSettings)
	g.GET("", a.handleGetSettings)
}

// handleSetSettings godoc
//
//	@Summary		Get the currently active Tenant Settings
//	@Security		apiKey
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=Settings}
//	@Router			/settings [get]
func (a *apiV1) handleGetSettings(c echo.Context) error {
	var tenant db.Tenant
	tenant.TenantID = db.ID(GetTenantContext(c).TenantID)
	res := a.db.Model(&db.Tenant{}).Find(&tenant)

	if res.Error != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, res.Error.Error()))
	}

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelope(c, tenant.TenantSettings))
}

// handleGetSettings godoc
//
//	@Summary		Apply new Tenant Settings
//	@Security		apiKey
//	@Param			settings body Settings true "New settings to apply"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=Settings}
//	@Router			/settings [post]
func (a *apiV1) handleSetSettings(c echo.Context) error {
	var settings Settings

	err := json.NewDecoder(c.Request().Body).Decode(&settings)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	var tenant db.Tenant
	tenant.TenantID = db.ID(GetTenantContext(c).TenantID)

	res := a.db.Model(&tenant).UpdateColumn("tenant_settings", &settings)

	if res.Error != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelope(c, "Updated Settings associated with the tenant"))
}
