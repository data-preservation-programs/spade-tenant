package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Wallets []Wallet

type Wallet struct {
	Address string `json:"address"`
	Active  bool   `json:"active"` // true - active dealmaking from this wallet, false - still counts as an associated wallet
}

// handleSetWallets godoc
//	@Summary		Update wallets used by a tenant
// 	@Param 		  token header string true "Auth token"
//  @Param 			wallets body Wallets true "New wallets to add or change Active flag of"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{Response=Wallets}
//	@Router			/wallet-addresses [put]
func (s *apiV1) handleSetWallets(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{})
}

// handleDeleteWallets godoc
//	@Summary		Delete wallets used by a tenant
// 	@Param 		  token header string true "Auth token"
//  @Param 			wallets body []string true "wallet addresses to delete"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{Response=Wallets}
//	@Router			/wallet-addresses [delete]
func (s *apiV1) handleDeleteWallets(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{})
}

// handleDeleteWallets godoc
//	@Summary		Get wallets used by a tenant
// 	@Param 		  token header string true "Auth token"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{Response=Wallets}
//	@Router			/wallet-addresses [get]
func (s *apiV1) handleGetWallets(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{})
}
