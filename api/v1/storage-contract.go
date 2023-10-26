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

type SetStorageContractResponse struct {
	Cid string `json:"cid"`
}

// TODO: Allow POST to set both the contract, OR the CID, OR if both are set -> verify match first (grab from IPFS) and serialize before doing anything

// handleSetStorageContract godoc
//	@Summary		Update storage contract
// 	@Description 	Updates the storage contract *note* this will require SPs to resubscribe if changed
// 	@Param 		  token header string true "Auth token"
//  @Param 			collection body StorageContract true "New Storage Contract to update to"
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
