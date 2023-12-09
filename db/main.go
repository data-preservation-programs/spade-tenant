package db

import (
	"github.com/data-preservation-programs/spade-tenant/config"
	gormigrate "github.com/go-gormigrate/gormigrate/v2"
	logging "github.com/ipfs/go-log/v2"
	"github.com/labstack/echo"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	log = logging.Logger("router")
	DB  = NewSpdTenantSvc().DB
)

type SpdTenantSvc struct {
	DB         *gorm.DB
	DryRunMode bool
	Config     config.TenantServiceConfig
}

func NewSpdTenantSvc() *SpdTenantSvc {
	config := config.InitConfig()
	if config.DEBUG {
		logging.SetDebugLogging()
	}

	dbi, err := OpenDatabase(config.DB_URL, config.DEBUG, config.DRY_RUN)
	if err != nil {
		log.Fatalf("could not open db: %s", err)
	}

	return &SpdTenantSvc{DB: dbi, DryRunMode: config.DRY_RUN, Config: config}
}

func SpdTenantSvcMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		c.Set("SERVICE_CONTEXT", NewSpdTenantSvc)

		return next(c)
	}
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
