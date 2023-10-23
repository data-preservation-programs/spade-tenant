package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Addresses []Address

type Address struct {
	Address   string `json:"address"`
	IsSigning bool   `json:"is_signing"` // true - active dealmaking from this address, false - still counts as an associated wallet
}

// handleSetAddress godoc
//	@Summary		Update addresses associated with a tenant
// 	@Param 		  token header string true "Auth token"
//  @Param 			addresses body Addresses true "New addresses to add or change is_signing flag of"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=Addresses}
//	@Router			/address [put]
func (s *apiV1) handleSetAddress(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{})
}

// handleDeleteAddress godoc
//	@Summary		Delete addresses used by a tenant
// 	@Param 		  token header string true "Auth token"
//  @Param 			addresses body []string true "addresses to delete"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=Addresses}
//	@Router			/address [delete]
func (s *apiV1) handleDeleteAddress(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{})
}

// handleGetAddress godoc
//	@Summary		Get addresses used by a tenant
// 	@Param 		  token header string true "Auth token"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=Addresses}
//	@Router			/address [get]
func (s *apiV1) handleGetAddress(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{})
}
