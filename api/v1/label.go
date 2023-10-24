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

// handleListConstraintLabel godoc
//	@Summary		List all constraint labels for the tenant
// 	@Param 		  token header string true "Auth token"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=ConstraintLabel}
//	@Router			/label/constraint/list [get]
func (s *apiV1) handleListConstraintLabel(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{})
}

// TODO: ValueLabel - separate APIs
