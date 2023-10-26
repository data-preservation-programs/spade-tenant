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

// handleSetEligibilityCriteria godoc
//	@Summary		Set sp eligibility criteria
// 	@Param 		  token header string true "Auth token"
//  @Param 			elibility_criteria body EligibilityCriteria true "New eligibility criteria to update to"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=EligibilityCriteria}
//	@Router			/eligibility-criteria [post]
func (s *apiV1) handleSetEligibilityCriteria(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{})
}

// handleGetEligibilityCriteria godoc
//	@Summary		Get sp eligibility criteria
// 	@Param 		  token header string true "Auth token"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=EligibilityCriteria}
//	@Router			/eligibility-criteria [get]
func (s *apiV1) handleGetEligibilityCriteria(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{})
}
