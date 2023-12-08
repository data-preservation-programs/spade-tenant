package db

import (
	"github.com/data-preservation-programs/spade-tenant/config"
	gormigrate "github.com/go-gormigrate/gormigrate/v2"
	logging "github.com/ipfs/go-log/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	log = logging.Logger("router")
	DB  *gorm.DB
)

var DDM *DeltaDM

func init() {
	DDM = NewDeltaDMDm()
}

type DeltaDM struct {
	DB         *gorm.DB
	DryRunMode bool
	Config     config.DeltaConfig
}

func NewDeltaDMDm() *DeltaDM {
	config := config.InitConfig()
	if config.Settings.DEBUG {
		logging.SetDebugLogging()
	}

	dbi, err := OpenDatabase(config.Settings.DB_URL, config.Settings.DEBUG, config.Settings.DRY_RUN)
	if err != nil {
		log.Fatalf("could not connect to delta api: %s", err)
	}

	return &DeltaDM{DB: dbi, DryRunMode: config.Settings.DRY_RUN, Config: config}
}

// Opens a database connection, and returns a gorm DB object.
func OpenDatabase(dbDsn string, debug bool, dryRun bool) (*gorm.DB, error) {
	var config = &gorm.Config{}
	if debug {
		config = &gorm.Config{
			// Logger: logger.Default.LogMode(logger.Info),
			DryRun: dryRun, // Don't apply to the db, just generate sql
		}
	}

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
