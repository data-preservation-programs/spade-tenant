package api

import (
	"math"
	"net/http"
	"time"

	"github.com/data-preservation-programs/spade-tenant/config"
	"github.com/labstack/echo/v4"
)

type apiV1 struct {
}

func NewApiV1() *apiV1 {
	return &apiV1{}
}

var FILECOIN_GENESIS_UNIX_EPOCH int64 = 1598306400

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

// @securityDefinitions.apiKey apiKey
// @type apiKey
// @in header
// @name Authorization
func (s *apiV1) RegisterRoutes(e *echo.Echo, config config.DeltaConfig) {
	e.GET("/status", s.handleStatus)

	e.Use(AuthMiddleware)

	ConfigureCollectionsRouter(e, s, config)
	ConfigureStorageContractRouter(e, s, config)
	ConfigureSPRouter(e, s, config)
	ConfigureAddressesRouter(e, s, config)
	ConfigureSettingsRouter(e, s, config)
	ConfigureConstraintLabelsRouter(e, s, config)
}

func ConfigureCollectionsRouter(e *echo.Echo, s *apiV1, config config.DeltaConfig) {
	e.POST("/collections", s.handleCreateCollection)
	e.GET("/collections", s.handleGetCollections)
	e.PUT("/collections", s.handleModifyCollection)
	e.DELETE("/collections", s.handleDeleteCollection)
}

func ConfigureStorageContractRouter(e *echo.Echo, s *apiV1, config config.DeltaConfig) {
	e.GET("/storage-contract", s.handleGetStorageContract)
	e.POST("/storage-contract", s.handleSetStorageContract)
}

func ConfigureSPRouter(e *echo.Echo, s *apiV1, config config.DeltaConfig) {
	e.GET("/sp", s.handleGetStorageProviders)
	e.POST("/sp", s.handleApproveStorageProviders)
	e.POST("/sp/suspend", s.handleSuspendStorageProviders)
	e.POST("/sp/unsuspend", s.handleUnsuspendStorageProvider)
	e.POST("/sp/eligibility-criteria", s.handleSetSpEligibilityCriteria)
	e.GET("/sp/eligibility-criteria", s.handleGetSpEligibilityCriteria)
}

func ConfigureAddressesRouter(e *echo.Echo, s *apiV1, config config.DeltaConfig) {
	e.PUT("/addresses", s.handleSetAddresses)
	e.DELETE("/addresses", s.handleDeleteAddresses)
	e.GET("/addresses", s.handleGetAddresses)
}

func ConfigureSettingsRouter(e *echo.Echo, s *apiV1, config config.DeltaConfig) {
	e.POST("/settings", s.handleGetSettings)
	e.GET("/settings", s.handleSetSettings)
}

func ConfigureConstraintLabelsRouter(e *echo.Echo, s *apiV1, config config.DeltaConfig) {
	e.GET("/constraint-labels", s.handleGetConstraintLabels)
}

func GetTenantId(c echo.Context) int {
	return int(c.Get(TENANT_CONTEXT).(AuthContext).TenantID)
}

func UnixToFilEpoch(unixEpoch int64) int64 {
	return int64(math.Floor(float64(unixEpoch-FILECOIN_GENESIS_UNIX_EPOCH) / 30))
}

func GetSlugFromErrorCode(errorCode int) string {
	switch errorCode {
	case 1:
		return "error_1"
	case 2:
		return "error_2"
	default:
		return "unknown"
	}
}

func CreateErrorResponseEnvelop(c echo.Context, errorCode int, err string) ResponseEnvelope {
	return ResponseEnvelope{
		RequestUUID:        c.Response().Header().Get(echo.HeaderXRequestID),
		ResponseTime:       time.Now(),
		ResponseStateEpoch: UnixToFilEpoch(time.Now().Unix()),
		ResponseCode:       http.StatusInternalServerError,
		ErrCode:            errorCode,
		ErrSlug:            GetSlugFromErrorCode(errorCode),
		Response:           nil,
	}
}

func CreateSuccessResponseEnvelop(c echo.Context, message interface{}) ResponseEnvelope {
	return ResponseEnvelope{
		RequestUUID:        c.Response().Header().Get(echo.HeaderXRequestID),
		ResponseTime:       time.Now(),
		ResponseStateEpoch: UnixToFilEpoch(time.Now().Unix()),
		ResponseCode:       http.StatusOK,
		Response:           message,
	}
}

type ResponseEnvelope struct {
	RequestUUID        string      `json:"request_uuid,omitempty"`
	ResponseTime       time.Time   `json:"response_timestamp"`
	ResponseStateEpoch int64       `json:"response_state_epoch,omitempty"`
	ResponseCode       int         `json:"response_code"`
	ErrCode            int         `json:"error_code,omitempty"`
	ErrSlug            string      `json:"error_slug,omitempty"`
	ErrLines           []string    `json:"error_lines,omitempty"`
	InfoLines          []string    `json:"info_lines,omitempty"`
	ResponseEntries    *int        `json:"response_entries,omitempty"`
	Response           interface{} `json:"response"`
}
