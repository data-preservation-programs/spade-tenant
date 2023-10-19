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

// handlePostCollections godoc
//	@Summary		Creates a new collection
// 	@Param 		  token header string true "Auth token"
//  @Param 			collection body CreateCollectionRequest true "Collection to create"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{Response=CreateCollectionResponse}
//	@Router			/collections [post]
func (s *apiV1) handleCreateCollection(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{})
}
