package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type MutableCollection struct {
	Name                   string                  `json:"name"`
	ReplicationConstraints []ReplicationConstraint `json:"replication_constraints"`
	PieceListSource        PieceListSource         `json:"piece_list_source"`
	Inactive               bool                    `json:"inactive"`
	DealParams             DealParams              `json:"deal_params"`
}

type DealParams struct {
	DurationDays     uint `json:"duration_days"`
	StartWithinHours uint `json:"start_within_hours"`
}

type Collection struct {
	UUID uuid.UUID `json:"uuid"`
	MutableCollection
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
	CollectionID uint   `json:"collection_id"`
	Status       string `json:"status"`
}

// handleCreateCollection godoc
//	@Summary		Creates a new collection
// 	@Param 		  token header string true "Auth token"
//  @Param 			collection body MutableCollection true "Collection to create"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=Collection}
//	@Router			/collections [post]
func (s *apiV1) handleCreateCollection(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{})
}

type GetCollectionsResponse []Collection

// handleGetCollections godoc
//	@Summary		Gets info about collections
// 	@Param 		  token header string true "Auth token"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=Collection}
//	@Router			/collections [get]
func (s *apiV1) handleGetCollections(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{})
}

// TODO: PUT to modify
// delete collection
// modify - change replication constraints, change name or change piece source
