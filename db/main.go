package db

import (
	gormigrate "github.com/go-gormigrate/gormigrate/v2"
	logging "github.com/ipfs/go-log/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	log = logging.Logger("router")
)

// Opens a database connection, and returns a gorm DB object.
func OpenDatabase(dbDsn string, debug bool) (*gorm.DB, error) {
	var config = &gorm.Config{}
	if debug {
		config = &gorm.Config{
			// Logger: logger.Default.LogMode(logger.Info),
			DryRun: true, // Don't apply to the db, just generate sql
		}
	}

	DB, err := gorm.Open(postgres.Open(dbDsn), config)

	m := gormigrate.New(DB, gormigrate.DefaultOptions, Migrations)

	// Initialization for fresh db (only run at first setup)
	m.InitSchema(BaselineSchema)

	// ! unique constraint, for each tenant/ValueID and tenant/label must be unique
	// DB.Table("labels").AddUniqueIndex("idx_labels_tenant_id_id", "tenant_id", "id")
	// DB.Table("labels").AddUniqueIndex("idx_labels_tenant_id_label", "tenant_id", "label")

	if err != nil {
		return nil, err
	}

	if err = m.Migrate(); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}
	log.Debugf("Migration ran successfully")

	return DB, nil
}
