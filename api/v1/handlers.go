package api

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type ResponseEnvelope struct {
	RequestID          string      `json:"request_id,omitempty"`
	ResponseTime       time.Time   `json:"response_timestamp"`
	ResponseStateEpoch int64       `json:"response_state_epoch,omitempty"`
	ResponseCode       int         `json:"response_code"`
	ErrCode            int         `json:"error_code,omitempty"`
	ErrSlug            string      `json:"error_slug,omitempty"`
	ErrLines           []string    `json:"error_lines,omitempty"`
	InfoLines          []string    `json:"info_lines,omitempty"`
	ResponseEntries    *int        `json:"response_entries,omitempty"`
	Response           interface{} `json:"response"`
}

type HealthResponse struct {
	Alive bool `json:"alive"`
}

// handleHealth godoc
//	@Summary		Simple health check endpoint
//	@Description	This endpoint is used to check the health of the service
// 	@Param 		  token header string true "Auth token"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{Response=HealthResponse}
//	@Failure		500	{object}	error
//	@Router			/handleHealth [get]
func (s *apiV1) handleHealth(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
	})
}
