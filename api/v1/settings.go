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

type Tenants struct {
	TenantID                 int          `json:"tenant_id"`
	TenantStorageContractCid string       `json:"tenant_storage_contract_cid"`
	TenantMeta               pgtype.JSONB `json:"tenant_meta"`
	TenantSettings           pgtype.JSONB `json:"tenant_settings"`
}

// handleSetSettings godoc
//
//	@Summary		Get the currently active Tenant Settings
//	@Param			token header string true "Auth token"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=Settings}
//	@Router			/settings [get]
func (s *apiV1) handleGetSettings(c echo.Context) error {
	var tenant Tenants
	res := db.DB.Where("tenant_id = ? ", GetTenantId(c)).Find(&tenant)
	if res.Error != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelop(c, res.Error.Error()))
	}

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelop(c, tenant.TenantSettings))
}

// handleGetSettings godoc
//
//	@Summary		Apply new Tenant Settings
//	@Param			token header string true "Auth token"
//	@Param 			settings body Settings true "New settings to apply"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=Settings}
//	@Router			/settings [post]
func (s *apiV1) handleSetSettings(c echo.Context) error {
	var settings pgtype.JSONB

	err := json.NewDecoder(c.Request().Body).Decode(&settings)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelop(c, err.Error()))
	}

	var tenant Tenants
	res := db.DB.Where("tenant_id = ? ", GetTenantId(c)).Find(&tenant)

	if res.Error != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelop(c, err.Error()))
	}

	if res.RowsAffected == 0 {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelop(c, "Tenant not found."))
	}

	res = db.DB.Model(&tenant).Where("tenant_id = ?", GetTenantId(c)).Update("tenant_settings", settings)

	if res.Error != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelop(c, err.Error()))
	}

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelop(c, "Updated Addresses associated with the tenant"))
}
