package utils

import (
	"testing"
)

func TestGetGetSubdivisionsByCountry(t *testing.T) {

	subdivisions, err := GetSubdivisionsByCountry("us")

	if err != nil || len(subdivisions) != 57 {
		t.Errorf("Expected USA to have 57 subdivisions")
	}
	subdivisions, err = GetSubdivisionsByCountry("usa")
	if err != nil || len(subdivisions) != 57 {
		t.Errorf("Expected USA to have 57 subdivisions")
	}

	subdivisions, err = GetSubdivisionsByCountry("USA")
	if err != nil || len(subdivisions) != 57 {
		t.Errorf("Expected USA to have 57 subdivisions")
	}

	subdivisions, err = GetSubdivisionsByCountry("United States")
	if err != nil || len(subdivisions) != 57 {
		t.Errorf("Expected USA to have 57 subdivisions")
	}

	subdivisions, err = GetSubdivisionsByCountry("United States of America")
	if err != nil || len(subdivisions) != 57 {
		t.Errorf("Expected USA to have 57 subdivisions")
	}

	_, err = GetSubdivisionsByCountry("Atlantis")
	if err != nil {
		t.Errorf("Expected Atlantis to not be valid")
	}
}

func TestGetCountryInfo(t *testing.T) {

	country := GetCountryInfo("us")
	if country == nil || country.Alpha2 != "US" {
		t.Errorf("Expected country to be US")
	}

	country = GetCountryInfo("usa")
	if country == nil || country.Alpha2 != "US" {
		t.Errorf("Expected country to be US")
	}

	country = GetCountryInfo("united states")
	if country == nil || country.Alpha2 != "US" {
		t.Errorf("Expected country to be US")
	}

	country = GetCountryInfo("canada")
	if country == nil || country.Alpha2 != "CA" {
		t.Errorf("Expected country to be CA")
	}

	country = GetCountryInfo("atlantis")
	if country != nil {
		t.Errorf("Expected country to not exist")
	}
}
