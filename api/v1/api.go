package api

import (
	"math"
	"net/http"
	"time"

	"github.com/data-preservation-programs/spade-tenant/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

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
func RegisterRoutes(e *echo.Echo, config config.TenantServiceConfig) {
	apiGroup := e.Group("/api/v1")
	e.Use(middleware.RequestID())
	e.Use(AuthMiddleware)

	ConfigureStatusRouter(apiGroup, config)
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

func CreateErrorResponseEnvelope(c echo.Context, errorCode int, err string) ResponseEnvelope {
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

func CreateSuccessResponseEnvelope(c echo.Context, message interface{}) ResponseEnvelope {
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
