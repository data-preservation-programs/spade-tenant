package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type SetPolicyRequest struct {
	Clauses []PolicyClause `json:"clauses"`
}

type PolicyClause struct {
	Attribute string      `json:"attribute"`
	Operator  string      `json:"operator"`
	Value     interface{} `json:"value"` // TODO
}

type SetPolicyResponse struct {
}

// handleSetPolicy godoc
//	@Summary		Set or update a policy
// 	@Param 		  token header string true "Auth token"
//  @Param 			collection body SetPolicyRequest true "New policy to update to"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{Response=SetPolicyResponse}
//	@Router			/policies [post]
func (s *apiV1) handleSetPolicy(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
	})
}

type SetStorageContractRequest struct {
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

// handleSetStorageContract godoc
//	@Summary		Update storage contract
// 	@Description 	Updates the storage contract *note* this will require SPs to resubscribe if changed
// 	@Param 		  token header string true "Auth token"
//  @Param 			collection body SetStorageContractRequest true "New Storage Contract to update to"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{Response=SetStorageContractResponse}
//	@Router			/policies/storage-contract [post]
func (s *apiV1) handleSetStorageContract(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
	})
}
