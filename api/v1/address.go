package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Addresses []Address
type AddressesMutable []AddressMutable

type Address struct {
	AddressMutable
	ActorID *uint `json:"actor_id,omitempty"` // TODO :swagger docs should be null not 0
}

type AddressMutable struct {
	Address   string `json:"address"`
	IsSigning bool   `json:"is_signing"` // true - active dealmaking from this address, false - still counts as an associated wallet
}

// handleSetAddresses godoc
//
//			@Summary		Update addresses associated with a tenant
//	 		@Security apiKey
//		  @Param 			addresses body AddressesMutable true "New addresses to add or change is_signing flag of"
//			@Produce		json
//			@Success		200	{object}	ResponseEnvelope{response=Addresses}
//			@Router			/addresses [put]
func (s *apiV1) handleSetAddresses(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{})
}

// handleDeleteAddresses godoc
//
//		@Summary		Delete addresses used by a tenant
//		@Security apiKey
//	  @Param 			addresses body []string true "addresses to delete"
//		@Produce		json
//		@Success		200	{object}	ResponseEnvelope{response=Addresses}
//		@Router			/addresses [delete]
func (s *apiV1) handleDeleteAddresses(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{})
}

// handleGetAddresses godoc
//
//	@Summary		Get addresses used by a tenant
//	@Security apiKey
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=Addresses}
//	@Router			/addresses [get]
func (s *apiV1) handleGetAddresses(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{})
}
