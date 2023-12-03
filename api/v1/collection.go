package api

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/data-preservation-programs/spade-tenant/db"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type WriteCollectionMutable struct {
	TenantID     int `json:"tenant_id"`
	CollectionID int `json:"collection_id"`
	CommonProperties
}

type CollectionMutable struct {
	CommonProperties
}

type CommonProperties struct {
	CollectionName        string          `json:"name"`
	CollectionPieceSource PieceListSource `json:"piece_list_source"`
	CollectionActive      bool            `json:"inactive"`
	CollectionDealParams  DealParams      `json:"deal_params"`
}

type DealParams struct {
	DurationDays     uint `json:"duration_days"`
	StartWithinHours uint `json:"start_within_hours"`
}

type Collection struct {
	UUID uuid.UUID `json:"uuid"`
	CollectionMutable
}

type PieceListSource struct {
	Method            string `json:"method"`
	PollIntervalHours int    `json:"poll_interval_hours,omitempty"`
	ConnectionDetails string `json:"connection_details"` // TODO: better types / validation for the connection details
}

type ReplicationConstraint struct {
	ConstraintID  int `json:"constraint_id"`
	ConstraintMax int `json:"constraint_max"`
}

type CreateCollectionResponse struct {
	CollectionID     uint   `json:"collection_id"`
	CollectionActive string `json:"status"`
}

// handleCreateCollection godoc
//
//		@Summary		Creates a new collection
//		@Param			token header string true "Auth token"
//	    @Param 			collection body MutableCollection true "Collection to create"
//		@Produce		json
//		@Success		200	{object}	ResponseEnvelope{response=Collection}
//		@Router			/collections [post]
func (s *apiV1) handleCreateCollection(c echo.Context) error {
	var collectionMutable CommonProperties
	err := json.NewDecoder(c.Request().Body).Decode(&collectionMutable)
	writeCollectionMutable := WriteCollectionMutable{TenantID: GetTenantId(c), CollectionID: rand.Intn(100), CommonProperties: collectionMutable}
	err = db.DB.Table("collections").Create(&writeCollectionMutable).Error

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelop(c, err.Error()))
	}

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelop(c, "Updated Addresses associated with the tenant"))
}

type GetCollectionsResponse []Collection

// handleGetCollections godoc
//
//	@Summary		Gets info about collections
//	@Param			token header string true "Auth token"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=Collection}
//	@Router			/collections [get]
func (s *apiV1) handleGetCollections(c echo.Context) error {
	// pieceListSource := PieceListSource{Method: "method", PollIntervalHours: 10, ConnectionDetails: "details"}

	// x, _ := json.Marshal(CommonProperties{CollectionName: "name", CollectionPieceSource: pieceListSource, CollectionActive: true, CollectionDealParams: DealParams{DurationDays: 10, StartWithinHours: 10}})
	// fmt.Println(string(x))
	var collectionResponse CreateCollectionResponse
	err := db.DB.Table("collections").Where("tenant_id = ? ", GetTenantId(c)).Find(&collectionResponse).Error

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelop(c, err.Error()))
	}

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelop(c, collectionResponse))
}

// handleModifyCollection godoc
//
//	@Summary		Modify a collection
//	@Param 			token header string true "Auth token"
//	@Param 			collectionUUID path string true "Collection UUID to modify"
//	@Param 			collection body MutableCollection true "Collection data to update"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=Collection}
//	@Router			/collections/:collectionUUID [put]
func (s *apiV1) handleModifyCollection(c echo.Context) error {
	fmt.Println("here")
	var collectionMutable CollectionMutable
	err := json.NewDecoder(c.Request().Body).Decode(&collectionMutable)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelop(c, err.Error()))
	}

	res := db.DB.Table("collections").Where("collection_id = ?", c.Param("collectionUUID")).Updates(&collectionMutable)

	if res.Error != nil {
		return res.Error
	}

	err = res.Save(&collectionMutable).Error

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelop(c, err.Error()))
	}

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelop(c, collectionMutable))
}

// handleDeleteCollection godoc
//
//	@Summary		Delete a collection
//	@Param 			token header string true "Auth token"
//	@Param 			collectionUUID path string true "Collection UUID to modify"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=bool}
//	@Router			/collections/:collectionUUID [delete]
func (s *apiV1) handleDeleteCollection(c echo.Context) error {
	var collectionMutable CollectionMutable
	err := json.NewDecoder(c.Request().Body).Decode(&collectionMutable)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelop(c, err.Error()))
	}

	err = db.DB.Table("collections").Where("collection_id = ?", c.Param("collectionUUID")).Delete(&WriteCollectionMutable{}).Error

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelop(c, err.Error()))
	}

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelop(c, "Deleted Collection associated with the tenant"))
}
