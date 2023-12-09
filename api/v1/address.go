package api

import (
	"encoding/json"
	"net/http"

	"github.com/data-preservation-programs/spade-tenant/db"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Addresses []db.Address
type Address db.Address

type AddressResponse struct {
	AddressRobust    string `json:"address_robust" gorm:"primaryKey"`
	AddressActorID   uint   `json:"actor_id"`
	AddressIsSigning bool   `json:"is_signing" gorm:"default:true;not null"`
}

var DB *gorm.DB

func ConfigureAddressesRouter(e *echo.Group, service *db.SpdTenantSvc) {
	DB = service.DB

	g := e.Group("/addresses")
	g.PUT("", handleSetAddresses)
	g.DELETE("", handleDeleteAddresses)
	g.GET("", handleGetAddresses)
}

// handleSetAddresses godoc
//
//	@Summary		Update addresses associated with a tenant
//	@Param			token header string true "Auth token"
//	@Param			addresses body AddressesMutable true "New addresses to add or change is_signing flag of"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=Addresses}
//	@Router			/addresses [put]
func handleSetAddresses(c echo.Context) error {
	var addresses Addresses
	err := json.NewDecoder(c.Request().Body).Decode(&addresses)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	err = DB.Transaction(func(tx *gorm.DB) error {
		for _, address := range addresses {
			address.TenantID = GetTenantId(c)
			res := tx.Save(&address)
			if res.Error != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelope(c, "Updated Addresses associated with the tenant"))
}

// handleDeleteAddresses godoc
//
//	@Summary		Delete addresses used by a tenant
//	@Param 			token header string true "Auth token"
//	@Param 			addresses body []string true "addresses to delete"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=Addresses}
//	@Router			/addresses [delete]
func handleDeleteAddresses(c echo.Context) error {
	var addressesIds []string
	err := json.NewDecoder(c.Request().Body).Decode(&addressesIds)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	err = DB.Where("tenant_id = ? AND address_robust in (?)", GetTenantId(c), addressesIds).Delete(&Address{}).Error

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelope(c, "Deleted Addresses associated with the tenant"))
}

// handleGetAddresses godoc
//
//	@Summary		Get addresses used by a tenant
//	@Param			token header string true "Auth token"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=Addresses}
//	@Router			/addresses [get]
func handleGetAddresses(c echo.Context) error {
	var addresses []AddressResponse
	res := DB.Table("addresses").Where("tenant_id = ? and deleted_at is NULL", GetTenantId(c)).Find(&addresses)

	if res.Error != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, res.Error.Error()))
	}

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelope(c, addresses))
}
