package api

import (
	"encoding/json"
	"net/http"

	"github.com/data-preservation-programs/spade-tenant/db"
	"github.com/jackc/pgtype"
	"github.com/labstack/echo/v4"
)

type Settings struct {
	SpAutoApprove  bool `json:"sp_auto_approve"`
	SpAutoSuspend  bool `json:"sp_auto_suspend"`
	MaxInFlightGiB uint `json:"max_in_flight_gib"`
}

func ConfigureSettingsRouter(e *echo.Group, service *db.SpdTenantSvc) {
	g := e.Group("/settings")
	g.POST("", handleSetSettings)
	g.GET("", handleGetSettings)
}

// handleSetSettings godoc
//
//	@Summary		Get the currently active Tenant Settings
//	@Security		apiKey
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=Settings}
//	@Router			/settings [get]
func handleGetSettings(c echo.Context) error {
	var tenant db.Tenant
	res := db.DB.Where("tenant_id = ? ", GetTenantId(c)).Find(&tenant)
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
func handleSetSettings(c echo.Context) error {
	var settings Settings

	err := json.NewDecoder(c.Request().Body).Decode(&settings)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	var tenant db.Tenant
	res := db.DB.Where("tenant_id = ? ", GetTenantId(c)).Find(&tenant)

	if res.Error != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	if res.RowsAffected == 0 {
		return c.JSON(http.StatusNotFound, CreateErrorResponseEnvelope(c, http.StatusNotFound, "Tenant not found."))
	}
	blob, err := json.Marshal(settings)

	tenant.TenantSettings = pgtype.JSONB{Bytes: blob}
	res = db.DB.Where("tenant_id = ?", GetTenantId(c)).Save(&tenant)

	if res.Error != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelope(c, "Updated Settings associated with the tenant"))
}
