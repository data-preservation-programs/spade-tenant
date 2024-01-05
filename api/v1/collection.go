package api

import (
	"encoding/json"
	"net/http"

	"github.com/data-preservation-programs/spade-tenant/db"
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/labstack/echo/v4"
)

type CollectionResponse struct {
	CollectionID          uuid.UUID    `json:"collection_id"`
	CollectionName        string       `json:"name"`
	CollectionActive      *bool        `json:"active"`
	CollectionPieceSource pgtype.JSONB `json:"piece_list_source"`
	CollectionDealParams  pgtype.JSONB `json:"deal_params"`
}

func (a *apiV1) ConfigureCollectionRouter(e *echo.Group) {
	g := e.Group("/collections")
	g.GET("", a.handleGetCollections)
	g.PUT("/:collectionUUID", a.handleModifyCollection)
	g.POST("", a.handleCreateCollection)
	g.DELETE("/:collectionUUID", a.handleDeleteCollection)
}

// handleCreateCollection godoc
//
//	@Summary		Creates a new collection
//	@Security		apiKey
//	@Param 			collection CollectionMutable body true "Collection to create"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=Collection}
//	@Router			/collections [post]
func (a *apiV1) handleCreateCollection(c echo.Context) error {
	var collection db.Collection

	err := json.NewDecoder(c.Request().Body).Decode(&collection)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	collection.TenantID = db.ID(GetTenantContext(c).TenantID)
	collection.CollectionID = uuid.New()
	collection.ReplicationConstraints = []db.ReplicationConstraint{{CollectionID: collection.CollectionID, ConstraintID: 0, ConstraintMax: 10}}
	active := true
	collection.CollectionActive = &active

	err = a.db.Create(&collection).Error

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelope(c,
		&CollectionResponse{CollectionID: collection.CollectionID,
			CollectionName:        *collection.CollectionName,
			CollectionActive:      collection.CollectionActive,
			CollectionPieceSource: collection.CollectionPieceSource,
			CollectionDealParams:  collection.CollectionDealParams}))
}

// handleGetCollections godoc
//
//	@Summary		Gets info about collections
//	@Security		apiKey
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=Collection}
//	@Router			/collections [get]
func (a *apiV1) handleGetCollections(c echo.Context) error {

	var collectionResponse []CollectionResponse

	err := a.db.Model(&db.Collection{TenantID: db.ID(GetTenantContext(c).TenantID)}).Find(&collectionResponse).Error

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelope(c, collectionResponse))
}

// handleModifyCollection godoc
//
//	@Summary		Modify a collection
//	@Security apiKey
//	@Param 		  collectionUUID path string true "Collection UUID to modify"
//	@Param 		  collection body MutableCollection true "Collection data to update"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=Collection}
//	@Router			/collections/:collectionUUID [put]
func (a *apiV1) handleModifyCollection(c echo.Context) error {
	var collection CollectionResponse
	err := json.NewDecoder(c.Request().Body).Decode(&collection)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	id, err := uuid.Parse(c.Param("collectionUUID"))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	err = a.db.Model(&db.Collection{TenantID: db.ID(GetTenantContext(c).TenantID), CollectionID: id}).Updates(&collection).Error

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelope(c,
		&CollectionResponse{CollectionID: collection.CollectionID,
			CollectionName:        collection.CollectionName,
			CollectionActive:      collection.CollectionActive,
			CollectionPieceSource: collection.CollectionPieceSource,
			CollectionDealParams:  collection.CollectionDealParams}))
}

// handleDeleteCollection godoc
//
//	@Summary		Delete a collection
//	@Security apiKey
//	@Param 		  collectionUUID path string true "Collection UUID to modify"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=bool}
//	@Router			/collections/:collectionUUID [delete]
func (a *apiV1) handleDeleteCollection(c echo.Context) error {
	err := a.db.Delete(&db.Collection{TenantID: db.ID(GetTenantContext(c).TenantID), CollectionID: uuid.MustParse(c.Param("collectionUUID"))}).Error

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelope(c, "Deleted Collection associated with the tenant with id "+c.Param("collectionUUID")))
}
