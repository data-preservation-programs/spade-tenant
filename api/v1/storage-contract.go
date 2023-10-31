package api

import (
	"net/http"

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

type SetStorageContractRequest struct {
	Cid             string          `json:"cid"`
	StorageContract StorageContract `json:"storage_contract"`
}

type SetStorageContractResponse struct {
	Cid string `json:"cid"`
}

// handleSetStorageContract godoc
//	@Summary		Update storage contract
// 	@Description 	Updates the storage contract. <br/>
//  @Description <br/> *Note* this will require SPs to resubscribe if changed.
//  @Description <br/> *Note* CID is optional, if specified, then `storage_contract` becomes optional.
//  @Description If both are specified, then we will validate that the CID matches the proposed storage contract and return an error if not.
//  @Description If only CID is specified, then we will fetch it and update the storage contract to it.
// 	@Param 		  token header string true "Auth token"
//  @Param 			collection body SetStorageContractRequest true "New Storage Contract to update to"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=StorageContract}
//	@Router			/storage-contract [post]
func (s *apiV1) handleSetStorageContract(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{})
}

type GetStorageContractResponse struct {
	Cid string `json:"cid"`
	StorageContract
}

// handleGetStorageContract godoc
//	@Summary		Get tenant storage contract
// 	@Param 		  token header string true "Auth token"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=GetStorageContractResponse}
//	@Router			/storage-contract [get]
func (s *apiV1) handleGetStorageContract(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{})
}
