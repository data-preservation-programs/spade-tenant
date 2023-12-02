package api

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

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
	e.Use(AuthMiddleware)
	// /collections
	e.POST("/collections", s.handleCreateCollection)
	e.GET("/collections", s.handleGetCollections)
	e.PUT("/collections", s.handleModifyCollection)
	e.DELETE("/collections", s.handleDeleteCollection)

	// /storage-contract
	e.GET("/storage-contract", s.handleGetStorageContract)
	e.POST("/storage-contract", s.handleSetStorageContract)

	// /sp
	e.GET("/sp", s.handleGetStorageProviders)
	e.POST("/sp", s.handleApproveStorageProviders)
	e.POST("/sp/suspend", s.handleSuspendStorageProviders)
	e.POST("/sp/unsuspend", s.handleUnsuspendStorageProvider)
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
func GetTenantId(c echo.Context) int {
	return int(c.Get(TENANT_CONTEXT).(AuthContext).TenantID)
}

func CreateErrorResponseEnvelop(c echo.Context, err string) ResponseEnvelope {
	return ResponseEnvelope{
		RequestUUID:        c.Response().Header().Get(echo.HeaderXRequestID),
		ResponseTime:       time.Now(),
		ResponseStateEpoch: time.Now().UTC().UnixMilli(),
		ResponseCode:       http.StatusInternalServerError,
		ErrCode:            http.StatusInternalServerError,
		ErrSlug:            err,
		Response:           err,
	}
}

func CreateSuccessResponseEnvelop(c echo.Context, message interface{}) ResponseEnvelope {
	return ResponseEnvelope{
		RequestUUID:        c.Response().Header().Get(echo.HeaderXRequestID),
		ResponseTime:       time.Now(),
		ResponseStateEpoch: time.Now().UTC().UnixMilli(),
		ResponseCode:       http.StatusOK,
		Response:           message,
	}
}
