package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/data-preservation-programs/spade-tenant/db"
	"github.com/jackc/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (a *apiV1) ConfigureBrokerRouter(e *echo.Group) {
	e.GET("/tenant-state", a.handleGetTenantsState)
	e.POST("/subscription-event", a.handlePostBrokerSubscriptionEvent)
}

type BrokerResponse []TenantBrokerPayload

type TenantBrokerPayload struct {
	TenantID                 db.ID             `json:"tenant_id"`
	TenantStorageContractCID string            `json:"tenant_storage_contract_cid"`
	TenantSettings           db.TenantSettings `json:"tenant_settings"`
	TenantAddresses          []db.Address      `json:"tenant_addresses"`
	CandidateSPs             []CandidateSP     `json:"candidate_sps"`
	Collections              []db.Collection   `json:"collections"`
}

type CandidateSP struct {
	SPID                    db.ID        `json:"sp_id"`
	ProviderTenantState     string       `json:"provider_tenant_state"`
	ProviderTenantStateInfo string       `json:"provider_tenant_state_info"`
	AttributeValues         map[int]int  `json:"attribute_values"`
	ProviderMetadata        pgtype.JSONB `json:"provider_metadata"` // TODO: struct for what we expect here (i.e, max_bytes_in_flight)
}

// handleGetTenantsState
// *Not included in Swagger docs
//
//	@Summary		List of all tenants in a JSON object to be consumed by the broker.
//	@Security Boker API Key
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=BrokerResponse}
//	Route:		/broker/tenant-state [get]
func (a *apiV1) handleGetTenantsState(c echo.Context) error {
	var br BrokerResponse

	// Raw DB query to get all tenants and their associated data
	var queryResult []db.Tenant
	a.db.Preload(clause.Associations).Preload("Collections.ReplicationConstraints").Preload("SPs.SPAttributes").Find(&queryResult)

	// Convert the query result to the BrokerResponse struct
	for _, tenant := range queryResult {
		var tenantSettings db.TenantSettings
		err := tenant.TenantSettings.AssignTo(&tenantSettings)
		if err != nil {
			// unable to unmarshal tenant settings json - fallaback to empty settings
			log.Error("unable to extract tenant settings from db" + err.Error())
		}

		var payload TenantBrokerPayload = TenantBrokerPayload{
			TenantID:                 tenant.TenantID,
			TenantStorageContractCID: tenant.TenantStorageContractCid,
			TenantSettings:           tenantSettings,
			TenantAddresses:          tenant.TenantAddresses,
			Collections:              tenant.Collections,
		}

		// Assemble CandidateSPs for this tenant
		for _, sp := range tenant.SPs {
			var candidateSP CandidateSP = CandidateSP{
				SPID:                    sp.SPID,
				ProviderTenantState:     string(sp.TenantSpState),
				ProviderMetadata:        sp.TenantSpsMeta,
				ProviderTenantStateInfo: sp.TenantSpStateInfo,
				AttributeValues:         make(map[int]int),
			}

			// Construct the attribute values for this SP
			for _, attribute := range sp.SPAttributes {
				// candidateSP.AttributeValues[strconv.Itoa(int(attribute.AttributeLabelID))] = strconv.Itoa(int(attribute.AttributeValueID))
				candidateSP.AttributeValues[int(attribute.AttributeLabelID)] = int(attribute.AttributeValueID)
			}

			payload.CandidateSPs = append(payload.CandidateSPs, candidateSP)
		}

		br = append(br, payload)
	}

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelope(c, br))
}

/** example payload
 * [
 *   { "spid" : 1 , tenantID: "5", "storage_contract_cid": "bafy1234", "authorization": "abcd" }
 * ]
 */
type BrokerSubscriptionEvents []SubscriptionEventPayload

type SubscriptionEventPayload struct {
	SPID               db.ID  `json:"sp_id"`
	TenantID           db.ID  `json:"tenant_id"`
	StorageContractCID string `json:"storage_contract_cid"`
	Authorization      string `json:"authorization"`
}

// handlePostBrokerSubscriptionEvent
// *Not included in Swagger docs
//
//	@Summary		Process subscription events from the Broker, updating tenant state
//	@Security Broker API Key
//	@Param 		  broker_subscription_events body BrokerSubscriptionEvents true "broker subscription events"
//	@Produce		json
//	Route:		200	{object}	ResponseEnvelope{response=BrokerSubscriptionEvents}
func (a *apiV1) handlePostBrokerSubscriptionEvent(c echo.Context) error {
	var subscriptionEvents BrokerSubscriptionEvents

	err := json.NewDecoder(c.Request().Body).Decode(&subscriptionEvents)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	err = a.db.Transaction(func(tx *gorm.DB) error {
		for _, se := range subscriptionEvents {

			tenant := db.Tenant{}
			res := tx.Model(&db.Tenant{}).Where("tenant_id = ?", se.TenantID).First(&tenant)
			if res.Error != nil {
				// Tenant ID does not exist - this should not happen and represents invalid input
				return fmt.Errorf("unable to find tenant with id " + strconv.Itoa(int(se.TenantID)))
			}

			var tenantSettings db.TenantSettings
			var nextState db.TenantSpState = db.TenantSpStatePending

			err := tenant.TenantSettings.AssignTo(&tenantSettings)
			if err != nil {
				// unable to unmarshal tenant settings json - fallback to Pending state
				log.Error("unable to extract tenant settings from db: " + err.Error())
			} else {
				if tenantSettings.SpAutoApprove {
					nextState = db.TenantSpStateActive
				} else {
					nextState = db.TenantSpStatePending
				}
			}

			meta := db.TenantsSpsMeta{
				Signature:         se.Authorization,
				SignedContractCID: se.StorageContractCID,
			}
			var pgMeta pgtype.JSONB
			pgMeta.Set(meta)

			var tsp db.TenantsSPs
			res = tx.Model(&db.TenantsSPs{}).Where("tenant_id = ? AND sp_id = ?", se.TenantID, se.SPID).First(&tsp)
			if res.Error != nil {
				// Tenant-SP mapping does not exist - represents invalid input
				return fmt.Errorf("unable to find tenant sp with tenant id " + strconv.Itoa(int(se.TenantID)) + " and sp id " + strconv.Itoa(int(se.SPID)))
			}

			tsp.TenantSpState = nextState
			tsp.TenantSpsMeta = pgMeta

			res = tx.Save(&tsp)
			if res.Error != nil {
				log.Error("unable to update tenant sp state " + res.Error.Error())
				return fmt.Errorf("unable to update tenant sp state " + res.Error.Error())
			}
		}

		return nil
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelope(c, subscriptionEvents))
}
