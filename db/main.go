package db

import (
	gormigrate "github.com/go-gormigrate/gormigrate/v2"
	logging "github.com/ipfs/go-log/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	log = logging.Logger("router")
)

// Opens a database connection, and returns a gorm DB object.
func OpenDatabase(dbDsn string, debug bool, dryRun bool) (*gorm.DB, error) {
	var config = &gorm.Config{}

	if debug {
		config = &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
			DryRun: dryRun,
		}
	}

	config.DryRun = dryRun

	DB, err := gorm.Open(postgres.Open(dbDsn), config)

	m := gormigrate.New(DB, gormigrate.DefaultOptions, Migrations)

	// Initialization for fresh db (only run at first setup)
	m.InitSchema(BaselineSchema)

	if err != nil {
		return nil, err
	}

	if err = m.Migrate(); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}

	log.Debugf("Migration ran successfully")

	return DB, nil
}
