package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type CreateCollectionRequest struct {
	Name                   string                  `json:"name"`
	ReplicationConstraints []ReplicationConstraint `json:"replication_constraints"`
	PieceSource            PieceSource             `json:"piece_source"`
}

type PieceSource struct {
	Method string `json:"method"`
	// TODO: Poll interval hours <optional>
	// TODO: Piece Source
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
//  @Param 			collection body CreateCollectionRequest true "Collection to create"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=CreateCollectionResponse}
//	@Router			/collections [post]
func (s *apiV1) handleCreateCollection(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{})
}

type GetCollectionsResponse []Collection

type Collection struct {
	ID                     string                  `json:"id"` // todo: UUIDs
	Name                   string                  `json:"name"`
	ReplicationConstraints []ReplicationConstraint `json:"replication_constraints"`
}

// handleGetCollections godoc
//	@Summary		Gets info about collections
// 	@Param 		  token header string true "Auth token"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=GetCollectionsResponse}
//	@Router			/collections [get]
func (s *apiV1) handleGetCollections(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{})
}

// TODO - hosted piece source
// Swagger docs for hosted piece source to make sure it makes sense
// Think about how it fits together coherently re: one single source of data for premium tenants
