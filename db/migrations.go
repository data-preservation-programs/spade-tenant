package db

import (
	gormigrate "github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// If this runs, it means the database is empty. No migrations will be applied on top of it, as this sets up the database from scratch so it starts out "up to date"
func BaselineSchema(tx *gorm.DB) error {
	log.Debugf("first run: initializing database schema")

	err := tx.Exec("CREATE TYPE tenant_sp_state AS ENUM ('eligible', 'pending', 'active', 'suspended');").Error

	if err != nil {
		log.Fatalf("error creating enum: %s", err)
	}

	err = tx.AutoMigrate(&Tenant{}, &Address{}, &TenantSPEligibilityClauses{}, &Collection{}, &Label{}, &SP{}, &TenantsSPs{}, &ReplicationConstraint{})
	if err != nil {
		log.Fatalf("error applying initial schema: %s", err)
	}

	return nil
}

var Migrations []*gormigrate.Migration = []*gormigrate.Migration{
	// {
	// 	ID: "00",
	// 	Migrate: func(tx *gorm.DB) error {
	// 	},
	// 	Rollback: func(tx *gorm.DB) error {
	// 		return errors.New("rollback not supported")
	// 	},
	// },
}
