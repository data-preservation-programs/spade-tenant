package api

import (
	"encoding/json"
	"net/http"

	"github.com/data-preservation-programs/spade-tenant/db"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (a *apiV1) ConfigureSpEligibilityCriteriaRouter(e *echo.Group) {
	g := e.Group("/sp")
	g.POST("/eligibility-criteria", a.handleSetSpEligibilityCriteria)
	g.GET("/eligibility-criteria", a.handleGetSpEligibilityCriteria)
	g.DELETE("/eligibility-criteria/attribute/:attribute", a.handleDeleteSpEligibilityCriteria)
}

type TenantSPEligibilityClausesResponse struct {
	ClauseAttribute string                `json:"attribute"`
	ClauseOperator  db.ComparisonOperator `json:"operator"`
	ClauseValue     string                `json:"value"`
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

func (a *apiV1) handleSetSpEligibilityCriteria(c echo.Context) error {
	var eligibilityCriteria []db.TenantSPEligibilityClauses
	err := json.NewDecoder(c.Request().Body).Decode(&eligibilityCriteria)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	err = a.db.Transaction(func(tx *gorm.DB) error {
		for _, eligibilityClause := range eligibilityCriteria {
			eligibilityClause.TenantID = db.ID(GetTenantContext(c).TenantID)
			res := tx.Save(&eligibilityClause)
			if res.Error != nil {
				return res.Error
			}
		}

		return nil
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelope(c, "Updated Eligibility clauses associated with the tenant"))
}

// handleGetSpEligibilityCriteria godoc
//
//	@Summary		Get sp eligibility criteria
//	@Security		apiKey
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=EligibilityCriteria}
//	@Router			/sp/eligibility-criteria [get]
func (a *apiV1) handleGetSpEligibilityCriteria(c echo.Context) error {
	var eligibilityCriteria []TenantSPEligibilityClausesResponse
	var clause db.TenantSPEligibilityClauses
	clause.TenantID = db.ID(GetTenantContext(c).TenantID)

	res := a.db.Model(&db.TenantSPEligibilityClauses{}).Find(&eligibilityCriteria)

	if res.Error != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, res.Error.Error()))
	}

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelope(c, eligibilityCriteria))
}

// handleDeleteSpEligibilityCriteria godoc
//
//	@Summary		Get sp eligibility criteria
//	@Security		apiKey
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=EligibilityCriteria}
//	@Router			/sp/eligibility-criteria/attribute/:attribute [delete]
func (a *apiV1) handleDeleteSpEligibilityCriteria(c echo.Context) error {
	res := a.db.Find(&db.TenantSPEligibilityClauses{TenantID: db.ID(GetTenantContext(c).TenantID), ClauseAttribute: c.Param("attribute")})

	if res.Error != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, res.Error.Error()))
	}

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelope(c, "Deleted attribute: "+c.Param("attribute")))
}
