package api

import (
	"net/http"

	"github.com/data-preservation-programs/spade-tenant/db"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm/clause"
)

func (a *apiV1) ConfigureBrokerRouter(e *echo.Group) {
	g := e.Group("/broker")
	g.GET("", a.handleGetTenantsInformation)
	g.POST("", a.handlePostNotifyTenantService)
}

type BrokerResponse []TenantBrokerPayload

type TenantBrokerPayload struct {
	TenantID                 db.ID  `json:"tenant_id"`
	TenantStorageContractCID string `json:"tenant_storage_contract_cid"`
	TenantSettings           struct {
		DealLengthDays int `json:"deal_length_days"`
	}
	TenantAddresses []db.Address   `json:"tenant_addresses"`
	CandidateSPs    []CandidateSPs `json:"candidate_sps"`
	Collections     []db.Collection
}

type CandidateSPs struct {
	SPID                db.ID             `json:"sp_id"`
	ProviderTenantState string            `json:"provider_tenant_state"`
	AttributeValues     map[string]string `json:"attribute_values"`
	ProviderMetadata    struct {
		MaxBytesInFlight int `json:"max_bytes_in_flight"`
	}
}

// handleGetTenantsInformation godoc
//
//	@Summary		List of all tenants in a JSON object to be consumed by the broker.
//	@Security apiKey
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=BrokerResponse}
//	@Router			/broker [get]
func (a *apiV1) handleGetTenantsInformation(c echo.Context) error {
	// var br BrokerResponse

	var queryResult []db.Tenant
	a.db.Preload(clause.Associations).Preload("Collections.ReplicationConstraints").Find(&queryResult)
	// a.db.Preload(clause.Associations).Preload("Collections").Preload("SPs").Preload("SP").Preload("TenantAddresses").Preload("Labels").Preload("Collections.ReplicationConstraints").Find(&queryResult)

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelope(c, queryResult))
}

// handlePostNotifyTenantService godoc
//
//	@Summary		Allows the broker to notify the tenant service.
//	@Security apiKey
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=LabelsResponse}
//	@Router			/broker [post]
func (a *apiV1) handlePostNotifyTenantService(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, CreateErrorResponseEnvelope(c, http.StatusNotImplemented, ""))
}
