package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/data-preservation-programs/spade-tenant/db"
	"github.com/jackc/pgtype"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type TenantsSPsResponse struct {
	SPID          db.ID            `json:"sp_id" gorm:"primaryKey;column:sp_id"`
	TenantSpState db.TenantSpState `gorm:"type:tenant_sp_state;column:tenant_sp_state;default:eligible;not null"`
	TenantSpsMeta pgtype.JSONB     `gorm:"type:jsonb;default:'{}';not null"`
}

func (a *apiV1) ConfigureSPRouter(e *echo.Group) {
	g := e.Group("/sp")
	g.GET("", a.handleGetStorageProviders)
	g.POST("/approve", a.handleApproveStorageProviders)
	g.POST("/suspend", a.handleSuspendStorageProviders)
	g.POST("/unsuspend", a.handleUnsuspendStorageProvider)
}

// @Summary		Get list of Storage Providers in all states
// @Security	apiKey
// @Produce		json
// @Success		200	{object}	ResponseEnvelope{response=GetStorageProvidersResponse}
// @Router		/sp [get]
func (a *apiV1) handleGetStorageProviders(c echo.Context) error {
	var storageProviderIds []TenantsSPsResponse

	a.db.Model(&db.TenantsSPs{TenantID: db.ID(GetTenantContext(c).TenantID)}).Find(&storageProviderIds)

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelope(c, storageProviderIds))
}

// handleApproveStorageProviders godoc
//
//	@Summary		Approves a list of Storage Providers to work with the tenant
//	@Description	Note: This is only required if auto_approve is false, requiring manual approval of SP subscription
//	@Security		apiKey
//	@Param			body StorageProviderIDs true "List of SP IDs to approve"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=GetStorageProvidersResponse}
//	@Router			/sp/approve [post]
func (a *apiV1) handleApproveStorageProviders(c echo.Context) error {
	storageProviderIds, err := a.updateTenantSpState(c, db.TenantSpStatePending, db.TenantSpStateActive)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelope(c, storageProviderIds))
}

// handleSuspendStorageProviders godoc
//
//	@Summary		Suspend storage providers
//	@Description	Note: This is only required if auto_suspend is false, as manual suspension is required
//	@Security		apiKey
//	@Param			body StorageProviderIDs true "List of SP IDs to suspend"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=GetStorageProvidersResponse}
//	@Router			/sp/suspend [post]
func (a *apiV1) handleSuspendStorageProviders(c echo.Context) error {
	storageProviderIds, err := a.updateTenantSpState(c, db.TenantSpStateActive, db.TenantSpStateSuspended)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelope(c, storageProviderIds))
}

// handleUnsuspendStorageProvider godoc
//
//	@Summary		Unsuspend a storage provider
//	@Security		apiKey
//	@Param			body StorageProviderIDs true "List of SP IDs to unsuspend"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=GetStorageProvidersResponse}
//	@Router			/sp/unsuspend [post]
func (a *apiV1) handleUnsuspendStorageProvider(c echo.Context) error {
	storageProviderIds, err := a.updateTenantSpState(c, db.TenantSpStateSuspended, db.TenantSpStateActive)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelope(c, storageProviderIds))
}

func (a *apiV1) updateTenantSpState(c echo.Context, stateToCheck db.TenantSpState, stateToSet db.TenantSpState) ([]int, error) {
	var storageProviderIds []int

	err := json.NewDecoder(c.Request().Body).Decode(&storageProviderIds)
	if err != nil {
		return storageProviderIds, err
	}
	err = a.db.Transaction(func(tx *gorm.DB) error {
		for _, id := range storageProviderIds {
			var tenantSps db.TenantsSPs
			a.db.Model(&db.TenantsSPs{SPID: db.ID(id), TenantID: db.ID(GetTenantContext(c).TenantID)}).Find(&tenantSps)
			if tenantSps.TenantSpState == stateToCheck {
				a.db.Model(&db.TenantsSPs{SPID: db.ID(id), TenantID: db.ID(GetTenantContext(c).TenantID)}).UpdateColumn("tenant_sp_state", stateToSet)
			} else {
				return fmt.Errorf("tenant sp state does not match. Expected {%s} but got {%s}", stateToCheck, tenantSps.TenantSpState)
			}

		}

		return nil
	})

	return storageProviderIds, err
}
