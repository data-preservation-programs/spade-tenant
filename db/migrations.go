package db

import (
	"fmt"
	"os"
	"strings"

	gormigrate "github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// If this runs, it means the database is empty. No migrations will be applied on top of it, as this sets up the database from scratch so it starts out "up to date"
func BaselineSchema(tx *gorm.DB) error {
	guardMigration("for new database/baseline")

	err := baselineSchemaManualMigrations(tx)
	if err != nil {
		log.Fatalf("error applying basline manual migrations: %s", err)
	}

	// Full table set-up
	err = tx.AutoMigrate(&Tenant{}, &Address{}, &TenantSPEligibilityClauses{}, &Collection{}, &Label{}, &SP{}, &TenantsSPs{}, &ReplicationConstraint{})
	if err != nil {
		log.Fatalf("error applying initial schema: %s", err)
	}

	return nil
}

// Baseline migrations that do not occur as part of `AutoMigrate`
func baselineSchemaManualMigrations(tx *gorm.DB) error {
	// Enums
	err := tx.Exec("CREATE TYPE tenant_sp_state AS ENUM ('eligible', 'pending', 'active', 'suspended');").Error
	if err != nil {
		return fmt.Errorf("error creating enum: %s", err)
	}
	err = tx.Exec("CREATE TYPE comparison_operator AS ENUM ('>', '<', '=', '>=', '<=', 'in', 'nin', '!=');").Error
	if err != nil {
		return fmt.Errorf("error creating enum: %s", err)
	}

	return nil
}

var Migrations []*gormigrate.Migration = []*gormigrate.Migration{
	// {
	// 	ID: "2023060800", // Set to todays date, starting with 00 for first migration
	// 	Migrate: func(tx *gorm.DB) error {
	// 		guardMigration("2023060800")
	// 		return tx.Migrator().AddColumn(&ReplicationConstraint{}, "RC")
	// 	},
	// },
}

// Ensure migrations are only run when env variable is specified
func guardMigration(migration string) {
	// TODO: grab from global config insead of directly accessing env
	enabled := strings.ToUpper(os.Getenv("DB_ALLOW_MIGRATIONS")) == "TRUE"

	if !enabled {
		log.Fatalf("Unable to apply migration %s as migrations are disabled. Set DB_ALLOW_MIGRATIONS=true to enable", migration)
	}
}
