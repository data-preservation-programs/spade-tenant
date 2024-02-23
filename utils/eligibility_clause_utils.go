package utils

import (
	"fmt"

	"github.com/biter777/countries"
)

func GetSubdivisionsByCountry(country string) ([]countries.SubdivisionCode, error) {
	if countries.AllSubdivisionsByCountryCode()[countries.ByName(country)][0] == "Unknown" {
		return []countries.SubdivisionCode{}, fmt.Errorf("country does not exist")
	}
	return countries.AllSubdivisionsByCountryCode()[countries.ByName(country)], nil
}

func GetAllCountryInfo() []*countries.Country {
	return countries.AllInfo()
}

func GetCountryInfo(country string) *countries.Country {
	for _, x := range countries.All() {
		if countries.ByName(country) == countries.ByName(x.Alpha2()) {
			return x.Info()
		}
	}

	return nil
}

type ClauseAttribute string

const (
	AgentCountry ClauseAttribute = "agentCountry"
	AgentRegion  ClauseAttribute = "agentRegion"
	AgentCity    ClauseAttribute = "agentCity"
)

func IsValidClauseAttribute(attribute string) bool {
	m := map[ClauseAttribute]bool{
		AgentCountry: true,
		AgentRegion:  true,
		AgentCity:    true,
	}

	_, valid := m[ClauseAttribute(attribute)]

	return valid
}
