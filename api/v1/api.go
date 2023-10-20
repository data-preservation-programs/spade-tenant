package api

import "github.com/labstack/echo/v4"

type apiV1 struct {
}

func NewApiV1() *apiV1 {
	return &apiV1{}
}

// @title Spade Tenant API
// @version 1.0.0
// @description This is the API for the Spade Tenant Application
// @termsOfService https://spade.storage

// @contact.name API Support
// @contact.url https://docs.spade.storage

// @license.name Apache 2.0 Apache-2.0 OR MIT
// @license.url https://github.com/data-preservation-programs/spade/blob/master/LICENSE.md

// @host api.spade.storage
// @BasePath  /
// @securityDefinitions.Bearer
// @securityDefinitions.Bearer.type apiKey
// @securityDefinitions.Bearer.in header
// @securityDefinitions.Bearer.name Authorization
func (s *apiV1) RegisterRoutes(e *echo.Echo) {
	e.GET("/health", s.handleHealth)

	// /collections
	e.POST("/collections", s.handleCreateCollection)
	e.GET("/collections", s.handleGetCollections)

	// /policy
	e.POST("/policy", s.handleSetPolicy)
	e.POST("/policy/storage-contract", s.handleSetStorageContract)
	e.GET("/policy", s.handleGetPolicy)
	e.GET("/policy/storage-contract", s.handleGetStorageContract)

	// /storage-providers
	e.GET("/storage-providers/subscribed", s.handleGetSubscribedStorageProviders)
	e.GET("/storage-providers/eligible", s.handleGetEligibleStorageProviders)
	e.POST("/storage-providers/approve", s.handleApproveStorageProviders)

	// /wallet-addresses
	e.PUT("/wallet-addresses", s.handleSetWallets)
	e.DELETE("/wallet-addresses", s.handleDeleteWallets)
	e.GET("/wallet-addresses", s.handleGetWallets)

	// /tenant
	e.POST("/tenant/settings", s.handleSetTenantSettings)
	e.GET("/tenant/settings", s.handleGetTenantSettings)

	// /mailbox
	e.GET("/mailbox", s.handleGetMailbox)
}
