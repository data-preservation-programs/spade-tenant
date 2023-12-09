package api

import (
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
)

var TENANT_CONTEXT = "TENANT_CONTEXT"

type AuthResponse struct {
	Result AuthResult `json:"result"`
}

type AuthResult struct {
	Validated bool `json:"validated"`
	Meta      struct {
		TenantID int32 `json:"tenant_id"`
	} `json:"details"`
}

type AuthContext struct {
	TenantID int32
}

// Secure routes by requiring a valid auth token
// This will place the Tenant's metadata into the context for use in downstream handlers
// To access them, call `tc := c.Get(TENANT_CONTEXT).(AuthContext)`
func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authKey, err := extractAuthKey(c.Request().Header.Get("Authorization"))

		if err != nil {
			return c.JSON(401, err.Error())
		}

		res, err := checkAuthToken(*authKey)
		if err != nil {
			return c.JSON(401, err.Error())
		}

		if !res.Validated || res.Meta.TenantID == 0 {
			return c.JSON(401, "invalid auth key")
		}

		authContext := AuthContext{
			TenantID: res.Meta.TenantID,
		}

		c.Set(TENANT_CONTEXT, authContext)

		return next(c)
	}
}

// Check that an auth string is populated in header and formatted correctly, then return it
//
//	`hint: pass in the value of c.Request().Header.Get("Authorization")`
func extractAuthKey(authorizationString string) (*string, error) {
	if authorizationString == "" {
		return nil, fmt.Errorf("missing auth header")
	}

	authParts := strings.Split(authorizationString, " ")
	if len(authParts) != 2 {
		return nil, fmt.Errorf("malformed auth header - must be of the form BEARER <token>")
	}
	if authParts[0] != "Bearer" {
		return nil, fmt.Errorf("malformed auth header - must have `Bearer` prefix")
	}

	return &authParts[1], nil
}

// Check the DB to see if a token is valid, return back TenantID
func checkAuthToken(token string) (*AuthResult, error) {
	// TODO: Query the DB for this Tenant's auth token and return
	res := AuthResult{
		Validated: true,
		Meta: struct {
			TenantID int32 `json:"tenant_id"`
		}{TenantID: 1},
	}

	return &res, nil
}
