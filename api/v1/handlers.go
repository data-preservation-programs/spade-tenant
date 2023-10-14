package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// handleHealth godoc
// @Summary      Simple health check endpoint
// @Description  This endpoint is used to check the health of the service
// @Produce      json
// @Success      200     {object}  string
// @Failure      500  {object} error
// @Router       /handleHealth [get]
func (s *apiV1) handleHealth(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
	})
}
