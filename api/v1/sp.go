package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/data-preservation-programs/spade-tenant/db"
	"github.com/labstack/echo/v4"
)

type StorageProvidersResponse []StorageProvider

type StorageProvider struct {
	SPID              uint      `json:"sp_id"`
	FirstActivatedAt  time.Time `json:"first_activated_at"`
	StatusLastChanged time.Time `json:"status_last_changed"`
	// State:
	// * eligible: The SP is eligible to work with the tenant, but has not yet begun the subscription process
	// * pending: The SP has begun the subscription process, but has not yet been approved by the tenant (note: only valid if auto-approve=false)
	// * active: The SP is active and working with the tenant. Deals may be made with this SP.
	// * suspended: The SP is suspended and deals may not be made with this SP, until it is un-suspended
	State string `json:"state" enums:"eligible,pending,active,suspended"`
}

func ConfigureSPRouter(e *echo.Group, service *db.SpdTenantSvc) {
	g := e.Group("/sp")
	g.GET("", handleGetStorageProviders)
	g.POST("", handleApproveStorageProviders)
	g.POST("/suspend", handleSuspendStorageProviders)
	g.POST("/unsuspend", handleUnsuspendStorageProvider)
}

// todo should work on tenant sp table
type StorageProviderIDs struct {
	SPIDs []uint `json:"sp_ids"`
}

// @Summary		Get list of Storage Providers
// @Security	apiKey
// @Produce		json
// @Success		200	{object}	ResponseEnvelope{response=GetStorageProvidersResponse}
// @Router		/sp [get]
func handleGetStorageProviders(c echo.Context) error {
	var storageProviderIds StorageProviderIDs

	err := json.NewDecoder(c.Request().Body).Decode(&storageProviderIds)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	var storageProviderResponse StorageProvidersResponse
	db.DB.Table("sps").Where("sp_id in (?)", storageProviderIds.SPIDs).Find(&storageProviderResponse)

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelope(c, storageProviderResponse))
}

// handleApproveStorageProviders godoc
//
//	@Summary		Approves a list of Storage Providers to work with the tenant
//	@Description Note: This is only required if auto_approve is false, requiring manual approval of SP subscription
//	@Security apiKey
//	@Param body body StorageProviderIDs true "List of SP IDs to approve"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=GetStorageProvidersResponse}
//	@Router			/sp/approve [post]
func handleApproveStorageProviders(c echo.Context) error {
	var storageProviderIds StorageProviderIDs

	err := json.NewDecoder(c.Request().Body).Decode(&storageProviderIds)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	var storageProviderResponse StorageProvidersResponse
	db.DB.Table("sps").Where("sp_id in (?)", storageProviderIds.SPIDs).UpdateColumn("tenant_sp_state", "eligible") //@jcace is active the right state

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelope(c, storageProviderResponse))
}

// handleSuspendStorageProviders godoc
//
//	@Summary		Suspend storage providers
//	@Description Note: This is only required if auto_suspend is false, as manual suspension is required
//	@Security apiKey
//	@Param body body StorageProviderIDs true "List of SP IDs to suspend"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=GetStorageProvidersResponse}
//	@Router			/sp/suspend [post]
func handleSuspendStorageProviders(c echo.Context) error {
	var storageProviderIds StorageProviderIDs

	err := json.NewDecoder(c.Request().Body).Decode(&storageProviderIds)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	var storageProviderResponse StorageProvidersResponse
	db.DB.Table("sps").Where("sp_id in (?)", storageProviderIds.SPIDs).UpdateColumn("tenant_sp_state", "suspended") //@jcace is active the right state

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelope(c, storageProviderResponse))
}

// handleUnsuspendStorageProvider godoc
//
//	@Summary		Unsuspend a storage provider
//	@Security apiKey
//	@Param body body StorageProviderIDs true "List of SP IDs to unsuspend"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=GetStorageProvidersResponse}
//	@Router			/sp/unsuspend [post]
func handleUnsuspendStorageProvider(c echo.Context) error {
	var storageProviderIds StorageProviderIDs

	err := json.NewDecoder(c.Request().Body).Decode(&storageProviderIds)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	var storageProviderResponse StorageProvidersResponse
	db.DB.Table("sps").Where("sp_id in (?)", storageProviderIds.SPIDs).UpdateColumn("tenant_sp_state", "active") //@jcace is active the right state

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelope(c, storageProviderResponse))
}
