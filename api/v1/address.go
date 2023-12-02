package api

import (
	"encoding/json"
	"net/http"

	"github.com/data-preservation-programs/spade-tenant/db"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Addresses []Address
type AddressesMutable []AddressMutable

type Address struct {
	AddressMutable
	TenantID       int
	AddressActorID *uint32 `json:"address_actor_id,omitempty"` // TODO :swagger docs should be null not 0
}

type AddressMutable struct {
	Address          string `json:"address"`
	AddressIsSigning bool   `json:"address_is_signing"` // true - active dealmaking from this address, false - still counts as an associated wallet
}

// handleSetAddresses godoc
//
//	@Summary		Update addresses associated with a tenant
//	@Param			token header string true "Auth token"
//	@Param			addresses body AddressesMutable true "New addresses to add or change is_signing flag of"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=Addresses}
//	@Router			/addresses [put]
func (s *apiV1) handleSetAddresses(c echo.Context) error {
	var addresses Addresses
	err := json.NewDecoder(c.Request().Body).Decode(&addresses)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelop(c, err.Error()))
	}

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		for _, address := range addresses {
			address.TenantID = GetTenantId(c)
			res := tx.Where("tenant_id = ? and address = ? ", address.TenantID, address.Address).Updates(&address)
			if res.Error != nil {
				return res.Error
			}

			if res.RowsAffected == 0 {
				tx.Create(&address)
			} else {
				res.Save(&address)
			}
		}

		return nil
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelop(c, err.Error()))
	}

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelop(c, "Updated Addresses associated with the tenant"))
}

// handleDeleteAddresses godoc
//
//	@Summary		Delete addresses used by a tenant
//	@Param 			token header string true "Auth token"
//	@Param 			addresses body []string true "addresses to delete"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=Addresses}
//	@Router			/addresses [delete]
func (s *apiV1) handleDeleteAddresses(c echo.Context) error {
	var addresses []string
	err := json.NewDecoder(c.Request().Body).Decode(&addresses)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelop(c, err.Error()))
	}

	err = db.DB.Where("tenant_id = ? AND address in (?)", GetTenantId(c), addresses).Delete(&Address{}).Error

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelop(c, err.Error()))
	}

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelop(c, "Deleted Addresses associated with the tenant"))
}

// handleGetAddresses godoc
//
//	@Summary		Get addresses used by a tenant
//	@Param			token header string true "Auth token"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=Addresses}
//	@Router			/addresses [get]
func (s *apiV1) handleGetAddresses(c echo.Context) error {
	var addresses Addresses
	res := db.DB.Where("tenant_id = ? ", GetTenantId(c)).Find(&addresses)

	if res.Error != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelop(c, res.Error.Error()))
	}

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelop(c, addresses))
}
