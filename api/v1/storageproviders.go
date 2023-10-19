package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type GetSubscribedStorageProvidersResponse struct {
	StorageProviderInfo []StorageProviderInfo `json:"storage_provider_info"`
}

type StorageProviderInfo struct {
	SPID              string   `json:"sp_id"`
	BytesStored       uint64   `json:"bytes_stored"`
	CidsStored        uint64   `json:"cids_stored"`
	CollectionsStored []string `json:"collections_stored"`
	SubscriptionDate  string   `json:"subscription_date"`
	LastDealDate      string   `json:"last_deal_date"`
	// TODO: SP Retrieval testing info (bswap, rrhttp, ipni)
}

// handleGetSubscribedStorageProviders godoc
//	@Summary		Get list of Storage Providers already working with the tenant, and their stats
// 	@Param 		  token header string true "Auth token"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{Response=GetSubscribedStorageProvidersResponse}
//	@Router			/storage-providers/subscribed [get]
func (s *apiV1) handleGetSubscribedStorageProviders(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{})
}

type GetEligibleStorageProvidersResponse struct {
	EligibleStorageProviders []EligibleStorageProvider `json:"eligible_storage_providers"`
}

type EligibleStorageProvider struct {
	SPID       string `json:"sp_id"`
	Subscribed bool   `json:"subscribed"`
}

// handleGetEligibleStorageProviders godoc
//	@Summary		Get list of Storage Providers not yet working with the tenant
// 	@Param 		  token header string true "Auth token"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{Response=GetEligibleStorageProvidersResponse}
//	@Router			/storage-providers/eligible [get]
func (s *apiV1) handleGetEligibleStorageProviders(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{})
}

type ApproveStorageProviders struct {
	SPIDs []string `json:"sp_ids"`
}

// handleApproveStorageProviders godoc
//	@Summary		Approves a list of Storage Providers to work with the tenant
// 	@Description Note: This is only required if auto_approve is false, requiring manual approval of SP subscription
// 	@Param 		  token header string true "Auth token"
//	@Param body body ApproveStorageProviders true "List of SP IDs to approve"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{Response=GetEligibleStorageProvidersResponse}
//	@Router			/storage-providers/approve [post]
func (s *apiV1) handleApproveStorageProviders(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{})
}
