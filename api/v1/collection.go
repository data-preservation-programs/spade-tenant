package api

import (
	"encoding/json"
	"net/http"

	"github.com/data-preservation-programs/spade-tenant/db"
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/labstack/echo/v4"
)

type CollectionMutable struct {
	CollectionName        string       `json:"name"`
	CollectionActive      bool         `json:"active"`
	CollectionPieceSource pgtype.JSONB `json:"piece_list_source"`
	CollectionDealParams  pgtype.JSONB `json:"deal_params"`

	ReplicationConstraints []db.ReplicationConstraint `json:"replication_constraints"`
}

// handleCreateCollection godoc
//
//	@Summary		Creates a new collection
//	@Security		apiKey
//	@Param 			collection CollectionMutable body true "Collection to create"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=Collection}
//	@Router			/collections [post]
func handleCreateCollection(c echo.Context) error {
	var collection db.Collection
	err := json.NewDecoder(c.Request().Body).Decode(&collection)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}
	collection.TenantID = GetTenantId(c)
	collection.CollectionID = uuid.New()

	err = db.DB.Create(&collection).Error

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelope(c, "Updated Addresses associated with the tenant"))
}

func ConfigureCollectionRouter(e *echo.Group, service *db.SpdTenantSvc) {
	g := e.Group("/collections")
	g.GET("", handleGetCollections)
	g.PUT("", handleModifyCollection)
	g.DELETE("", handleDeleteCollection)
}

// handleGetCollections godoc
//
//	@Summary		Gets info about collections
//	@Security		apiKey
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=Collection}
//	@Router			/collections [get]
func handleGetCollections(c echo.Context) error {
	var collectionResponse []db.Collection
	err := db.DB.Omit("tenant_id").Where("tenant_id = ? ", GetTenantId(c)).Find(&collectionResponse).Error

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
func handleModifyCollection(c echo.Context) error {
	var collection db.Collection
	err := json.NewDecoder(c.Request().Body).Decode(&collection)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	err = db.DB.Where("collection_id = ?", c.Param("collectionUUID")).Updates(&collection).Error

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelope(c, collection))
}

// handleDeleteCollection godoc
//
//	@Summary		Delete a collection
//	@Security apiKey
//	@Param 		  collectionUUID path string true "Collection UUID to modify"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=bool}
//	@Router			/collections/:collectionUUID [delete]
func handleDeleteCollection(c echo.Context) error {
	var collection db.Collection
	err := json.NewDecoder(c.Request().Body).Decode(&collection)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	err = db.DB.Where("collection_id = ?", c.Param("collectionUUID")).Delete(&db.Collection{}).Error

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelope(c, "Deleted Collection associated with the tenant"))
}
