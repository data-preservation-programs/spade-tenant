package api

import (
	"encoding/json"
	"net/http"

	"github.com/data-preservation-programs/spade-tenant/db"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type EligibilityCriteria struct {
	Clauses []EligibilityClause `json:"clauses"`
}

type EligibilityClause struct {
	Attribute string      `json:"attribute"`
	Operator  string      `json:"operator"`
	TenantID  db.ID       `json:"tenant_id"`
	Value     interface{} `json:"value"` // TODO: type - either []string or string
}

type ReadEligibilityCriteria struct {
	Clauses []ReadEligibilityClause `json:"clauses"`
}

type ReadEligibilityClause struct {
	Attribute string      `json:"attribute"`
	Operator  string      `json:"operator"`
	Value     interface{} `json:"value"` // TODO: type - either []string or string
}

func ConfigureSpEligibilityCriteriaRouter(e *echo.Group, service *db.SpdTenantSvc) {
	g := e.Group("/sp")
	g.POST("/eligibility-criteria", handleSetSpEligibilityCriteria)
	g.GET("/eligibility-criteria", handleGetSpEligibilityCriteria)
}

// todo maybe an update 1-1
// handleSetSpEligibilityCriteria godoc
//
//	@Summary		Set sp eligibility criteria
//	@Security		apiKey
//	@Param			eligibility_criteria body EligibilityCriteria true "New eligibility criteria to update to"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=EligibilityCriteria}
//	@Router			/sp/eligibility-criteria [post]
//
// todo put 1-1
func handleSetSpEligibilityCriteria(c echo.Context) error {
	var eligibilityCriteria EligibilityCriteria
	err := json.NewDecoder(c.Request().Body).Decode(&eligibilityCriteria.Clauses)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		for _, eligibilityClause := range eligibilityCriteria.Clauses {
			eligibilityClause.TenantID = GetTenantId(c)
			res := tx.Table("tenant_sp_eligibility_clauses").Where(&eligibilityClause).Save(&eligibilityClause)
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
func handleGetSpEligibilityCriteria(c echo.Context) error {
	var eligibilityCriteria ReadEligibilityCriteria
	res := db.DB.Table("tenant_sp_eligibility_clauses").Where("tenant_id = ? ", GetTenantId(c)).Find(&eligibilityCriteria.Clauses)

	if res.Error != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, res.Error.Error()))
	}

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelope(c, eligibilityCriteria))
}
