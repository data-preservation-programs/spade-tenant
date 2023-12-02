package api

import (
	"net/http"

	"github.com/data-preservation-programs/spade-tenant/db"
	"github.com/labstack/echo/v4"
	"gorm.io/datatypes"
)

type ConstraintLabels []Label

type Label struct {
	TenantId  int            `json:"tenant_id"`
	LabelID   int            `json:"label_id"`
	LabelText string         `json:"label_text"`
	Options   datatypes.JSON `json:"options"` // example: {"CA": 10, "US": 20}
}

// handleGetConstraintLabels godoc
//
//	@Summary		List all constraint labels for the tenant
//	@Param 			token header string true "Auth token"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=ConstraintLabels}
//	@Router			/constraint-labels [get]
func (s *apiV1) handleGetConstraintLabels(c echo.Context) error {
	var constraintLabels ConstraintLabels
	res := db.DB.Where("tenant_id = ? ", GetTenantId(c)).Find(&constraintLabels)

	if res.Error != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelop(c, res.Error.Error()))
	}

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelop(c, constraintLabels))
}

// TODO: PUT - for V2

/** example response
[
  {
    "label": "country",
    "id": 1,
    "enum": {
      "CANADA": 10,
      "USA": 20
    }
  },
  {
    "label": "org", // does not need enum. just ints used for the matching
    "id": 2
  }
]
*/
