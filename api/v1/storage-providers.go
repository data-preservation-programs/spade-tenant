package api

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type GetStorageProvidersResponse []StorageProvider

type StorageProvider struct {
	SPID              uint      `json:"sp_id"`
	FirstActivatedAt  time.Time `json:"first_activated_at"`
	StatusLastChanged time.Time `json:"status_last_changed"`
	Status            string    `json:"status"` // * ENUM: [ eligible, pending-approval, active, suspended ]
}

// TODO: Put enum in swagger description

//	@Summary		Get list of Storage Providers
// 	@Param 		  token header string true "Auth token"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=GetStorageProvidersResponse}
//	@Router			/sp [get]
func (s *apiV1) handleGetStorageProviders(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{})
}

type StorageProviderIDs struct {
	SPIDs []uint `json:"sp_ids"`
}

// handleApproveStorageProviders godoc
//	@Summary		Approves a list of Storage Providers to work with the tenant
// 	@Description Note: This is only required if auto_approve is false, requiring manual approval of SP subscription
// 	@Param 		  token header string true "Auth token"
//	@Param body body StorageProviderIDs true "List of SP IDs to approve"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=GetStorageProvidersResponse}
//	@Router			/sp/approve [post]
func (s *apiV1) handleApproveStorageProviders(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{})
}

// handleSuspendStorageProviders godoc
//	@Summary		Suspend storage providers
// 	@Description Note: This is only required if auto_suspend is false, as manual suspension is required
// 	@Param 		  token header string true "Auth token"
//	@Param body body StorageProviderIDs true "List of SP IDs to suspend"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=GetStorageProvidersResponse}
//	@Router			/sp/suspend [post]
func (s *apiV1) handleSuspendStorageProviders(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{})
}
