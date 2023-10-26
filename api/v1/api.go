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
// @BasePath  /tenant
// @securityDefinitions.Bearer
// @securityDefinitions.Bearer.type apiKey
// @securityDefinitions.Bearer.in header
// @securityDefinitions.Bearer.name Authorization
func (s *apiV1) RegisterRoutes(e *echo.Echo) {
	e.GET("/status", s.handleStatus)

	// /collections
	e.POST("/collections", s.handleCreateCollection)
	e.GET("/collections", s.handleGetCollections)

	// /eligibility-criteria
	e.POST("/policy", s.handleSetEligibilityCriteria)
	e.GET("/policy", s.handleGetEligibilityCriteria)

	// /storage-contract
	e.GET("/storage-contract", s.handleGetStorageContract)
	e.POST("/storage-contract", s.handleSetStorageContract)

	// /storage-providers
	e.GET("/storage-providers", s.handleGetStorageProviders)
	e.POST("/storage-providers/approve", s.handleApproveStorageProviders)
	e.POST("/storage-providers/suspend", s.handleSuspendStorageProviders)

	// /address
	e.PUT("/address", s.handleSetAddress)
	e.DELETE("/address", s.handleDeleteAddress)
	e.GET("/address", s.handleGetAddress)

	// /tenant
	e.POST("/settings", s.handleGetSettings)
	e.GET("/settings", s.handleSetSettings)

	// /label
	e.GET("/label/constraint/list", s.handleGetConstraintLabels)
}
