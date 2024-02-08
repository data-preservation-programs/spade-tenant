package api

import (
	"encoding/json"
	"net/http"

	"github.com/data-preservation-programs/spade-tenant/db"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Addresses []db.Address
type Address db.Address

type AddressResponse struct {
	AddressRobust    string `json:"address_robust"`
	AddressActorID   uint   `json:"actor_id"`
	AddressIsSigning bool   `json:"is_signing"`
}

func (a *apiV1) ConfigureAddressesRouter(e *echo.Group) {
	g := e.Group("/addresses")
	g.PUT("", a.handleUpdateAddresses)
	g.POST("", a.handleCreateAddresses)
	g.DELETE("", a.handleDeleteAddresses)
	g.GET("", a.handleGetAddresses)
}

// handleSetAddresses godoc
//
//	@Summary		Update addresses associated with a tenant
//	@Security apiKey
//	@Param			addresses body Addresses true "New addresses to add or change is_signing flag of"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=Addresses}
//	@Router			/addresses [put]
func (a *apiV1) handleUpdateAddresses(c echo.Context) error {
	var addresses Addresses
	err := json.NewDecoder(c.Request().Body).Decode(&addresses)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	var response []string
	err = a.db.Transaction(func(tx *gorm.DB) error {
		for _, address := range addresses {
			address.TenantID = db.ID(GetTenantContext(c).TenantID)

			res := tx.Model(&address).Updates(&address)

			if res.Error != nil {
				return err
			}

			if res.RowsAffected > 0 {
				response = append(response, *address.AddressRobust)
			}
		}

		return nil
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelope(c, response))
}

// handleSetAddresses godoc
//
//	@Summary		Creates addresses associated with a tenant
//	@Security apiKey
//	@Param			addresses body AddressMutable true "New addresses to add"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=Addresses}
//	@Router			/addresses [post]
func (a *apiV1) handleCreateAddresses(c echo.Context) error {
	var addresses Addresses
	err := json.NewDecoder(c.Request().Body).Decode(&addresses)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	var response []string
	err = a.db.Transaction(func(tx *gorm.DB) error {
		for _, address := range addresses {
			address.TenantID = db.ID(GetTenantContext(c).TenantID)

			res := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&address)

			if res.Error != nil {
				return err
			}

			if res.RowsAffected > 0 {
				response = append(response, *address.AddressRobust)
			}
		}
		return nil
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelope(c, response))
}

// handleDeleteAddresses godoc
//
//	@Summary		Delete addresses used by a tenant
//	@Security apiKey
//	@Param 			addresses body []string true "addresses to delete"
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=Addresses}
//	@Router			/addresses [delete]
func (a *apiV1) handleDeleteAddresses(c echo.Context) error {
	var addressesIds []string
	err := json.NewDecoder(c.Request().Body).Decode(&addressesIds)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}
	var response []string
	err = a.db.Transaction(func(tx *gorm.DB) error {
		for _, address := range addressesIds {
			res := a.db.Delete(&Address{TenantID: db.ID(GetTenantContext(c).TenantID), AddressRobust: &address})

			if res.RowsAffected > 0 {
				response = append(response, address)
			}

			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelope(c, response))
}

// handleGetAddresses godoc
//
//	@Summary		Get addresses used by a tenant
//	@Security apiKey
//	@Produce		json
//	@Success		200	{object}	ResponseEnvelope{response=Addresses}
//	@Router			/addresses [get]
func (a *apiV1) handleGetAddresses(c echo.Context) error {
	var addresses []AddressResponse
	res := a.db.Model(Address{TenantID: db.ID(GetTenantContext(c).TenantID)}).Find(&addresses)

	if res.Error != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorResponseEnvelope(c, http.StatusInternalServerError, res.Error.Error()))
	}

	return c.JSON(http.StatusOK, CreateSuccessResponseEnvelope(c, addresses))
}
