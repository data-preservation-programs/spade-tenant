package db

import (
	"fmt"

	gormigrate "github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// If this runs, it means the database is empty. No migrations will be applied on top of it, as this sets up the database from scratch so it starts out "up to date"
func BaselineSchema(tx *gorm.DB) error {
	confirmMigration("for brand new database - Baseline")

	// Enums
	err := tx.Exec("CREATE TYPE tenant_sp_state AS ENUM ('eligible', 'pending', 'active', 'suspended');").Error
	if err != nil {
		log.Fatalf("error creating enum: %s", err)
	}

	err = tx.Exec("CREATE TYPE comparison_operator AS ENUM ('>', '<', '=', '>=', '<=', 'in', 'nin', '!=');").Error
	if err != nil {
		log.Fatalf("error creating enum: %s", err)
	}

	// Full table set-up
	err = tx.AutoMigrate(&Tenant{}, &Address{}, &TenantSPEligibilityClauses{}, &Collection{}, &Label{}, &SP{}, &TenantsSPs{}, &ReplicationConstraint{})
	if err != nil {
		log.Fatalf("error applying initial schema: %s", err)
	}

	return nil
}

var Migrations []*gormigrate.Migration = []*gormigrate.Migration{
	// {
	// 	ID: "2023060800", // Set to todays date, starting with 00 for first migration
	// 	Migrate: func(tx *gorm.DB) error {
	// 		confirmMigration("2023060800")
	// 		return tx.Migrator().AddColumn(&ReplicationConstraint{}, "RC")
	// 	},
	// },
}

// Confirm migrations
func confirmMigration(migrationName string) {
	fmt.Printf("Migration %s must be applied. Enter 'Y' to run migration or any other key to abort: \n", migrationName)
	var input string
	fmt.Scanln(&input)
	if input != "Y" {
		log.Fatal("user aborted")
	}
}
