package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// TODO: More formal mailbox message structure
type MailboxMessage struct {
	Notifications []interface{} `json:"notifications"`
}

// handleGetMailbox godoc
//	@Summary		Gets mailbox messages for the tenant
// 	@Param 		  token header string true "Auth token"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{Response=MailboxMessage}
//	@Router			/mailbox [get]
func (s *apiV1) handleGetMailbox(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
	})
}
