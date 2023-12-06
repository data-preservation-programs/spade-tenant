package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// handleStatus godoc
//
//	@Summary		Simple health check endpoint
//	@Description	This endpoint is used to check the health of the service
//	@Security apiKey
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=string}
//	@Router			/status [get]
func (s *apiV1) handleStatus(c echo.Context) error {
	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelop(c, "Healthy"))
}
