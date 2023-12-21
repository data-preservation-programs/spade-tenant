package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/data-preservation-programs/spade-tenant/db"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type GetStorageProvidersResponse []StorageProvider

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

// @Summary		Get list of Storage Providers
// @Security apiKey
// @Produce		json
// @Success		200	{object}	ResponseEnvelope{response=GetStorageProvidersResponse}
// @Router			/sp [get]
func handleGetStorageProviders(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{})
}

type StorageProviderIDs struct {
	SPIDs []uint `json:"sp_ids"`
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
func (a *apiV1) handleApproveStorageProviders(c echo.Context) error {
	var storageProviderIds []int

	err := json.NewDecoder(c.Request().Body).Decode(&storageProviderIds)
	if err != nil {
		return err
	}

	err = a.db.Transaction(func(tx *gorm.DB) error {
		for _, id := range storageProviderIds {
			err = a.db.Model(&db.TenantsSPs{SPID: db.ID(id), TenantID: db.ID(GetTenantContext(c).TenantID), TenantSpState: db.TenantSpStateActive}).Update("tenant_sp_state", db.TenantSpStateActive).Error

			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelope(c, storageProviderIds))
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
func (a *apiV1) handleSuspendStorageProviders(c echo.Context) error {
	var storageProviderIds []int

	err := json.NewDecoder(c.Request().Body).Decode(&storageProviderIds)
	if err != nil {
		return err
	}

	err = a.db.Transaction(func(tx *gorm.DB) error {
		for _, id := range storageProviderIds {
			err = a.db.Model(&db.TenantsSPs{SPID: db.ID(id), TenantID: db.ID(GetTenantContext(c).TenantID), TenantSpState: db.TenantSpStateSuspended}).Update("tenant_sp_state", db.TenantSpStateSuspended).Error

			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelope(c, storageProviderIds))
}

// handleUnsuspendStorageProvider godoc
//
//	@Summary		Unsuspend a storage provider
//	@Security apiKey
//	@Param body body StorageProviderIDs true "List of SP IDs to unsuspend"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=GetStorageProvidersResponse}
//	@Router			/sp/unsuspend [post]
func (a *apiV1) handleUnsuspendStorageProvider(c echo.Context) error {
	var storageProviderIds []int

	err := json.NewDecoder(c.Request().Body).Decode(&storageProviderIds)
	if err != nil {
		return err
	}

	err = a.db.Transaction(func(tx *gorm.DB) error {
		for _, id := range storageProviderIds {
			err = a.db.Model(&db.TenantsSPs{SPID: db.ID(id), TenantID: db.ID(GetTenantContext(c).TenantID), TenantSpState: db.TenantSpStateActive}).Update("tenant_sp_state", db.TenantSpStateActive).Error

			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelope(c, storageProviderIds))
}
