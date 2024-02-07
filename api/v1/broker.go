package api

import (
	"net/http"

	"github.com/data-preservation-programs/spade-tenant/db"
	"github.com/jackc/pgtype"
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
	TenantID                 db.ID         `json:"tenant_id"`
	TenantStorageContractCID string        `json:"tenant_storage_contract_cid"`
	TenantSettings           pgtype.JSONB  `json:"tenant_settings"` // TODO: struct for what we expect here (i.e, max_deal_length)
	TenantAddresses          []db.Address  `json:"tenant_addresses"`
	CandidateSPs             []CandidateSP `json:"candidate_sps"`
	Collections              []db.Collection
}

type CandidateSP struct {
	SPID                    db.ID             `json:"sp_id"`
	ProviderTenantState     string            `json:"provider_tenant_state"`
	ProviderTenantStateInfo string            `json:"provider_tenant_state_info"`
	AttributeValues         map[string]string `json:"attribute_values"`
	ProviderMetadata        pgtype.JSONB      `json:"provider_metadata"` // TODO: struct for what we expect here (i.e, max_bytes_in_flight)
}

// handleGetTenantsInformation godoc
//
//	@Summary		List of all tenants in a JSON object to be consumed by the broker.
//	@Security apiKey
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=BrokerResponse}
//	@Router			/broker [get]
func (a *apiV1) handleGetTenantsInformation(c echo.Context) error {
	var br BrokerResponse

	// Raw DB query to get all tenants and their associated data
	var queryResult []db.Tenant
	a.db.Preload(clause.Associations).Preload("Collections.ReplicationConstraints").Find(&queryResult)

	// Convert the query result to the BrokerResponse struct
	for _, tenant := range queryResult {
		var payload TenantBrokerPayload = TenantBrokerPayload{
			TenantID:                 tenant.TenantID,
			TenantStorageContractCID: tenant.TenantStorageContractCid,
			TenantSettings:           tenant.TenantSettings,
			TenantAddresses:          tenant.TenantAddresses,
			Collections:              tenant.Collections,
		}

		// Get sp relation table metadata for this tenant
		// Note: It appears that Gorm doesn't like having metadata attached to a many2many relation table, thus there's no clean way to access it from the previous query
		// For simplicity, we simply query this table directly for the Tenant, and use that to construct the CandidateSPs struct
		// It should be possible to be done in one query but will require some raw SQL
		var tenantSPsRelation []db.TenantsSPs
		a.db.Model(&db.TenantsSPs{SPID: tenant.TenantID}).Find(&tenantSPsRelation)

		// Construct the CandidateSPs struct for this tenant
		for _, sp := range tenant.SPs {
			var providerTenantState db.TenantSpState
			var providerTenantMeta pgtype.JSONB
			var providerTenantStateInfo string
			for _, tenantSP := range tenantSPsRelation {
				if sp.SPID == tenantSP.SPID {
					providerTenantState = tenantSP.TenantSpState
					providerTenantMeta = tenantSP.TenantSpsMeta
					providerTenantStateInfo = tenantSP.TenantSpStateInfo
				}
			}

			var candidateSP CandidateSP = CandidateSP{
				SPID:                    sp.SPID,
				ProviderTenantState:     string(providerTenantState),
				ProviderMetadata:        providerTenantMeta,
				ProviderTenantStateInfo: providerTenantStateInfo,
				// TODO: Attributes
			}
			payload.CandidateSPs = append(payload.CandidateSPs, candidateSP)
		}

		br = append(br, payload)
	}

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelope(c, br))
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
