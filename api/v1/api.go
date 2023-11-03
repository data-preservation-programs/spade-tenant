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

	// /storage-contract
	e.GET("/storage-contract", s.handleGetStorageContract)
	e.POST("/storage-contract", s.handleSetStorageContract)

	// /sp
	e.GET("/sp", s.handleGetStorageProviders)
	e.POST("/sp", s.handleApproveStorageProviders)
	e.POST("/sp/suspend", s.handleSuspendStorageProviders)
	e.POST("/sp/eligibility-criteria", s.handleSetSpEligibilityCriteria)
	e.GET("/sp/eligibility-criteria", s.handleGetSpEligibilityCriteria)

	// /addresses
	e.PUT("/addresses", s.handleSetAddresses)
	e.DELETE("/addresses", s.handleDeleteAddresses)
	e.GET("/addresses", s.handleGetAddresses)

	// /settings
	e.POST("/settings", s.handleGetSettings)
	e.GET("/settings", s.handleSetSettings)

	// /constraint-labels
	e.GET("/constraint-labels", s.handleGetConstraintLabels)
}
