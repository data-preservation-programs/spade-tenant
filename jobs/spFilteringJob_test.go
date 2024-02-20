package main

import (
	"testing"

	"github.com/data-preservation-programs/spade-tenant/db"
)

func TestCompare(t *testing.T) {
	us := []string{"US"}

	if !compare(us, "US", db.IncludedIn) {
		t.Errorf("Expected %q to contain US", us)
	}

	if compare(us, "US", db.ExcludedFrom) {
		t.Errorf("Expected %q to be excluded from contain US", us)
	}

	countries := []string{"AH", "CA", "US", "DE", "MX"}

	if !compare(countries, "CA", db.IncludedIn) {
		t.Errorf("Expected %q to contain CA", us)
	}

	if !compare(countries, "MC", db.ExcludedFrom) {
		t.Errorf("Expected %q to be excluded from contain MC", us)
	}
}

func TestIsValidClause(t *testing.T) {
	clauses := []db.TenantSPEligibilityClauses{{TenantID: 1, ClauseAttribute: "AgentCountry", ClauseOperator: "in", ClauseValue: []string{"US"}}}
	usSP := SP{Provider: "f123", AgentCountry: "US"}
	if !verifyEligibility(clauses, usSP) {
		t.Errorf("Expected %s to be a valid provider", usSP.Provider)
	}

	caSP := SP{Provider: "f0343", AgentCountry: "CA", AgentCity: "Mariposa"}
	if verifyEligibility(clauses, caSP) {
		t.Errorf("Expected %s to be a valid provider", caSP.Provider)
	}

	clauses = []db.TenantSPEligibilityClauses{{TenantID: 1, ClauseAttribute: "AgentCountry", ClauseOperator: "nin", ClauseValue: []string{"US"}}}
	if !verifyEligibility(clauses, caSP) {
		t.Errorf("Expected %s to be a valid provider", caSP.Provider)
	}

	clauses = []db.TenantSPEligibilityClauses{{TenantID: 1, ClauseAttribute: "AgentCountry", ClauseOperator: "nin", ClauseValue: []string{"US"}},
		{TenantID: 1, ClauseAttribute: "AgentCity", ClauseOperator: "in", ClauseValue: []string{"Mariposa"}}}
	if !verifyEligibility(clauses, caSP) {
		t.Errorf("Expected %s to be a valid provider", caSP.Provider)
	}
}

func TestCreateTenantClauseMap(t *testing.T) {
	clauses := []db.TenantSPEligibilityClauses{{TenantID: 1, ClauseAttribute: "AgentCountry", ClauseOperator: "nin", ClauseValue: []string{"US"}},
		{TenantID: 1, ClauseAttribute: "AgentCity", ClauseOperator: "in", ClauseValue: []string{"Mariposa"}},
		{TenantID: 2, ClauseAttribute: "AgentCity", ClauseOperator: "in", ClauseValue: []string{"Chilly Beach"}}}
	res := createTenantClauseMap(clauses)

	if len(res) != 2 {
		t.Errorf("tenantClauseMap should be length 2 but got %d", len(res))
	}
}
