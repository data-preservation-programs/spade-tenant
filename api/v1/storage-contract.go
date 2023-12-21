package api

import (
	"encoding/json"
	"net/http"

	"github.com/data-preservation-programs/spade-tenant/db"
	"github.com/ipfs/go-cid"
	"github.com/labstack/echo/v4"
)

type StorageContract struct {
	Content struct {
		InfoLines []string `json:"info_lines"`
	}
	Retrieval struct {
		Mechanisms struct {
			IpldBitswap bool `json:"ipld_bitswap"`
			PieceRrhttp bool `json:"piece_rrhttp"`
		}
		Sla struct {
			InfoLines []string `json:"info_lines"`
		}
	}
}

type AddressedStorageContract struct {
	Cid             string          `json:"cid"`
	StorageContract StorageContract `json:"storage_contract"`
}

func (a *apiV1) ConfigureStorageContractRouter(e *echo.Group) {
	g := e.Group("/storage-contract")
	g.POST("", a.handleSetStorageContract)
	g.GET("", a.handleGetStorageContract)
}

// handleSetStorageContract godoc
//
//	@Summary		Update storage contract
//	@Description 	Updates the storage contract. <br/>
//	@Description <br/> *Note* this will require SPs to resubscribe if changed.
//	@Description <br/> *Note* CID is optional, if specified, then `storage_contract` becomes optional.
//	@Description If both are specified, then we will validate that the CID matches the proposed storage contract and return an error if not.
//	@Description If only CID is specified, then we will fetch it and update the storage contract to it.
//	@Security		apiKey
//	@Param 			collection body AddressedStorageContract true "New Storage Contract to update to"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=AddressedStorageContract}
//	@Router			/storage-contract [post]
func (a *apiV1) handleSetStorageContract(c echo.Context) error {
	var addressedStorageContract AddressedStorageContract
	err := json.NewDecoder(c.Request().Body).Decode(&addressedStorageContract)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	_, err = cid.Decode(addressedStorageContract.Cid)

	if err != nil || len(addressedStorageContract.Cid) == 0 {
		return c.JSON(http.StatusNotImplemented, CreateErrorResponseEnvelope(c, http.StatusNotImplemented, "StorageContract is not currently supported. Please pass in a CID."))
	}

	var tenant db.Tenant
	tenant.TenantID = db.ID(GetTenantContext(c).TenantID)
	err = a.db.Find(&tenant).Error

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	tenant.TenantStorageContractCid = addressedStorageContract.Cid
	err = a.db.Model(&tenant).Updates(&tenant).Error

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelope(c, "Set Storage Contract CID to "+addressedStorageContract.Cid))
}

// handleGetStorageContract godoc
//
//	@Summary		Get tenant storage contract
//	@Security		apiKey
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=AddressedStorageContract}
//	@Router			/storage-contract [get]
func (a *apiV1) handleGetStorageContract(c echo.Context) error {
	var tenant db.Tenant
	tenant.TenantID = db.ID(GetTenantContext(c).TenantID)
	err := a.db.Find(&tenant).Error

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelope(c, tenant.TenantStorageContractCid))
}
