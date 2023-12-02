package api

import (
	"encoding/json"
	"net/http"

	"github.com/data-preservation-programs/spade-tenant/db"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type WriteEligibilityCriteria struct {
	Clauses []WriteEligibilityClause `json:"clauses"`
}

type WriteEligibilityClause struct {
	TenantID int `json:"tenant_id"`
	CommonAttributes
}

type ReadEligibilityCriteria struct {
	Clauses []ReadEligibilityClause `json:"clauses"`
}

type ReadEligibilityClause struct {
	CommonAttributes
}

type CommonAttributes struct {
	ClauseAttribute string `json:"attribute"`
	ClauseOperator  string `json:"operator"`
	ClauseValue     string `json:"value"` // TODO: type - either []string or string
}

// handleSetSpEligibilityCriteria godoc
//
//	@Summary		Set sp eligibility criteria
//	@Param			token header string true "Auth token"
//	@Param			elibility_criteria body EligibilityCriteria true "New eligibility criteria to update to"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=EligibilityCriteria}
//	@Router			/sp/eligibility-criteria [post]
func (s *apiV1) handleSetSpEligibilityCriteria(c echo.Context) error {

	var eligibilityCriteria WriteEligibilityCriteria
	err := json.NewDecoder(c.Request().Body).Decode(&eligibilityCriteria.Clauses)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelop(c, err.Error()))
	}

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		for _, eligibilityClause := range eligibilityCriteria.Clauses {
			eligibilityClause.TenantID = GetTenantId(c)
			res := tx.Table("tenant_sp_eligibility_clauses").Where(&eligibilityClause).Updates(&eligibilityClause)
			if res.Error != nil {
				return res.Error
			}

			if res.RowsAffected == 0 {
				tx.Table("tenant_sp_eligibility_clauses").Create(&eligibilityClause)
			} else {
				res.Table("tenant_sp_eligibility_clauses").Save(&eligibilityClause)
			}
		}

		return nil
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelop(c, err.Error()))
	}

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelop(c, "Updated Eligibility clauses associated with the tenant"))
}

// handleGetSpEligibilityCriteria godoc
//
//	@Summary		Get sp eligibility criteria
//	@Param 			token header string true "Auth token"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=EligibilityCriteria}
//	@Router			/sp/eligibility-criteria [get]
func (s *apiV1) handleGetSpEligibilityCriteria(c echo.Context) error {
	var eligibilityCriteria ReadEligibilityCriteria
	res := db.DB.Table("tenant_sp_eligibility_clauses").Where("tenant_id = ? ", GetTenantId(c)).Find(&eligibilityCriteria.Clauses)

	if res.Error != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelop(c, res.Error.Error()))
	}

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelop(c, eligibilityCriteria))
}
