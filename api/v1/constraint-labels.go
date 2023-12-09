package api

import (
	"net/http"

	"github.com/data-preservation-programs/spade-tenant/db"
	"github.com/labstack/echo/v4"
)

type ConstraintLabels []Label

type Label struct {
	LabelId      int             `json:"label_id"`
	LabelText    string          `json:"labe_text"`
	LabelOptions map[string]uint `json:"labe_options"` // example: {"CA": 10, "US": 20}
}

func ConfigureSpConstraintLabelsRouter(e *echo.Group, service *db.SpdTenantSvc) {
	g := e.Group("/constraint-labels")
	g.GET("", handleGetConstraintLabels)
}

// handleGetConstraintLabels godoc
//
//	@Summary		List all constraint labels for the tenant
//	@Security apiKey
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=ConstraintLabels}
//	@Router			/constraint-labels [get]
func handleGetConstraintLabels(c echo.Context) error {
	var constraintLabels ConstraintLabels
	res := db.DB.Table("labels").Where("tenant_id = ? ", GetTenantId(c)).Find(&constraintLabels)

	if res.Error != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, res.Error.Error()))
	}

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelope(c, constraintLabels))
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
