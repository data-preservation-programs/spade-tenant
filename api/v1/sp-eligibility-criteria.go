package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type EligibilityCriteria struct {
	Clauses []EligibilityClause `json:"clauses"`
}

type EligibilityClause struct {
	Attribute string `json:"attribute"`
	Operator  string `json:"operator"`

	Value interface{} `json:"value"` // TODO: type - either []string or string
}

// handleSetSpEligibilityCriteria godoc
//
//	@Summary		Set sp eligibility criteria
//	@Security		apiKey
//	@Param			eligibility_criteria body EligibilityCriteria true "New eligibility criteria to update to"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=EligibilityCriteria}
//	@Router			/sp/eligibility-criteria [post]
//

func handleSetSpEligibilityCriteria(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{})
}

// handleGetSpEligibilityCriteria godoc
//
//	@Summary		Get sp eligibility criteria
//	@Security apiKey
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=EligibilityCriteria}
//	@Router			/sp/eligibility-criteria [get]
func handleGetSpEligibilityCriteria(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{})
}
