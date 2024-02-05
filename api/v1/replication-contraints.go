package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/data-preservation-programs/spade-tenant/db"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ReplicationConstraints []ReplicationConstraint
type ReplicationConstraint db.ReplicationConstraint

type ReplicationConstraintResponse struct {
	CollectionID  uuid.UUID `json:"collection_id"`
	ConstraintID  db.ID     `json:"constraint_id"`
	ConstraintMax int       `json:"constraint_max"`
}

func (a *apiV1) ConfigureReplicationConstraintsRouter(e *echo.Group) {
	g := e.Group("/collections/:collectionUUID/replication-constraints")
	g.PUT("", a.handleSetReplicationConstraints)
	g.DELETE("", a.handleDeleteReplicationConstraints)
	g.GET("", a.handleGetReplicationConstraints)
}

// handleSetReplicationConstraints godoc
//
//	@Summary		Creates or updates ReplicationConstraints associated with a tenant and collection
//	@Param			token header string true "Auth token"
//	@Param			replication_constraints body ReplicationConstraints true "New replication constraints to add or change"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=ReplicationConstraints}
//	@Router			/collections/:collectionUUID/replication-constraints [put]
func (a *apiV1) handleSetReplicationConstraints(c echo.Context) error {
	var replicationConstraints []ReplicationConstraintResponse

	err := json.NewDecoder(c.Request().Body).Decode(&replicationConstraints)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	id, err := uuid.Parse(c.Param("collectionUUID"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	err = a.db.Transaction(func(tx *gorm.DB) error {
		for i, replicationConstraint := range replicationConstraints {
			replicationConstraints[i].CollectionID = id
			constraint := &db.ReplicationConstraint{TenantID: db.ID(GetTenantContext(c).TenantID), CollectionID: id, ConstraintID: replicationConstraint.ConstraintID, ConstraintMax: replicationConstraint.ConstraintMax}

			res := tx.Save(&constraint)

			if res.Error != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelope(c, replicationConstraints))
}

// handleDeleteReplicationConstraints godoc
//
//	@Summary		Deletes ReplicationConstraints used by a tenant.
//	@Param			token header string true "Auth token"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=ReplicationConstraints}
//	@Router			/collections/:collectionUUID/replication-constraints [delete]
func (a *apiV1) handleDeleteReplicationConstraints(c echo.Context) error {
	var ids []int

	err := json.NewDecoder(c.Request().Body).Decode(&ids)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	id, err := uuid.Parse(c.Param("collectionUUID"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	var response []string
	err = a.db.Transaction(func(tx *gorm.DB) error {
		for _, constraintId := range ids {
			if constraintId > 0 {
				res := a.db.Delete(&ReplicationConstraint{TenantID: db.ID(GetTenantContext(c).TenantID), CollectionID: id, ConstraintID: db.ID(constraintId)})

				if res.RowsAffected > 0 {
					response = append(response, fmt.Sprint(constraintId))
				}
			}

			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelope(c, response))
}

// handleGetReplicationConstraints godoc
//
//	@Summary		Get ReplicationConstraints used by a tenant
//	@Param			token header string true "Auth token"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=ReplicationConstraints}
//	@Router			/collections/:collectionUUID/replication-constraints [get]
func (a *apiV1) handleGetReplicationConstraints(c echo.Context) error {
	var replicationConstraints []ReplicationConstraintResponse

	id, err := uuid.Parse(c.Param("collectionUUID"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	res := a.db.Model(&db.ReplicationConstraint{TenantID: db.ID(GetTenantContext(c).TenantID), CollectionID: id}).Find(&replicationConstraints)

	if res.Error != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, res.Error.Error()))
	}

	if res.RowsAffected == 0 {
		return c.JSON(http.StatusNotFound, CreateErrorResponseEnvelope(c, http.StatusNotFound, "Replication constraint not found"))
	}

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelope(c, replicationConstraints))
}
