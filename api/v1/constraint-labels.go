package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ConstraintLabels []Label

type Label struct {
	UUID    uuid.UUID       `json:"uuid"`
	Label   string          `json:"label"`
	Options map[string]uint `json:"options"` // example: {"CA": 10, "US": 20}
}

// handleGetConstraintLabels godoc
//
//	@Summary		List all constraint labels for the tenant
//	@Param 			token header string true "Auth token"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=ConstraintLabels}
//	@Router			/constraint-labels [get]
func (s *apiV1) handleGetConstraintLabels(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{})
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
