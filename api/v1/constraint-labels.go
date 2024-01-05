package api

import (
	"net/http"

	"github.com/data-preservation-programs/spade-tenant/db"
	"github.com/jackc/pgtype"
	"github.com/labstack/echo/v4"
)

type LabelResponse struct {
	LabelID      db.ID  `json:"id"`
	LabelText    string `json:"label"`
	LabelOptions pgtype.JSONB
}

func (a *apiV1) ConfigureSpConstraintLabelsRouter(e *echo.Group) {
	g := e.Group("/constraint-labels")
	g.GET("", a.handleGetConstraintLabels)
}

// handleGetConstraintLabels godoc
//
//	@Summary		List all constraint labels for the tenant
//	@Security apiKey
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=ConstraintLabels}
//	@Router			/constraint-labels [get]
func (a *apiV1) handleGetConstraintLabels(c echo.Context) error {
	var constraintLabels LabelResponse

	res := a.db.Model(&db.Label{TenantID: db.ID(GetTenantContext(c).TenantID)}).Find(&constraintLabels)

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
