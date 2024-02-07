package api

import (
	"net/http"
	"time"

	"github.com/data-preservation-programs/spade-tenant/config"
	"github.com/data-preservation-programs/spade-tenant/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

type apiV1 struct {
	db  *gorm.DB
	cfg *config.TenantServiceConfig
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

// @securityDefinitions.apiKey apiKey
// @type apiKey
// @in header
// @name Authorization
func (a *apiV1) RegisterRoutes(e *echo.Echo) {
	tenantApiGroup := e.Group("/api/v1")
	brokerApiGroup := e.Group("/api/v1/broker")
	e.Use(middleware.RequestID())

	// Separate group, as it needs separate auth middleware
	brokerApiGroup.Use(BrokerAuthMiddleware(a.cfg.BROKER_ACCESS_TOKEN))
	a.ConfigureBrokerRouter(brokerApiGroup)

	tenantApiGroup.Use(AuthMiddleware)
	ConfigureStatusRouter(tenantApiGroup)
	a.ConfigureAddressesRouter(tenantApiGroup)
	a.ConfigureStorageContractRouter(tenantApiGroup)
	a.ConfigureSPRouter(tenantApiGroup)
	a.ConfigureSettingsRouter(tenantApiGroup)
	a.ConfigureSpEligibilityCriteriaRouter(tenantApiGroup)
	a.ConfigureSpConstraintLabelsRouter(tenantApiGroup)
	a.ConfigureCollectionRouter(tenantApiGroup)
	a.ConfigureReplicationConstraintsRouter(tenantApiGroup)

}

func NewApiV1(db *gorm.DB, cfg *config.TenantServiceConfig) *apiV1 {
	return &apiV1{db: db, cfg: cfg}
}

func GetTenantContext(c echo.Context) AuthContext {
	return c.Get(TENANT_CONTEXT).(AuthContext)
}

func GetSlugFromErrorCode(errorCode int) string {
	//TODO Agree on what the slug to error codes we want to use are.
	switch errorCode {
	case 1:
		return "error_1"
	case 2:
		return "error_2"
	default:
		return "unknown"
	}
}

func CreateErrorResponseEnvelope(c echo.Context, errorCode int, err string) ResponseEnvelope {
	return ResponseEnvelope{
		RequestUUID:        c.Response().Header().Get(echo.HeaderXRequestID),
		ResponseTime:       time.Now(),
		ResponseStateEpoch: utils.UnixToFilEpoch(time.Now().Unix()),
		ResponseCode:       errorCode,
		ErrCode:            errorCode,
		ErrSlug:            GetSlugFromErrorCode(errorCode),
		Response:           nil,
	}
}

func CreateSuccessResponseEnvelope(c echo.Context, message interface{}) ResponseEnvelope {
	return ResponseEnvelope{
		RequestUUID:        c.Response().Header().Get(echo.HeaderXRequestID),
		ResponseTime:       time.Now(),
		ResponseStateEpoch: utils.UnixToFilEpoch(time.Now().Unix()),
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
