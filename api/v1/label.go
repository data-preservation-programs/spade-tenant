package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ConstraintLabel []Label

type Label struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

// handleGetConstraintLabels godoc
//	@Summary		List all constraint labels for the tenant
// 	@Param 		  token header string true "Auth token"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=ConstraintLabel}
//	@Router			/constraint-labels [get]
func (s *apiV1) handleGetConstraintLabels(c echo.Context) error {
	// TODO:return both constraint type, and value labels from one endpoint
	// nested json
	return c.JSON(http.StatusNotImplemented, map[string]string{})
}

// TODO: PUT

/**
 * // common ones only
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
