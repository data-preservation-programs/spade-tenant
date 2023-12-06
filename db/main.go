package db

import (
	"os"
	"strconv"

	"github.com/data-preservation-programs/spade-tenant/initializers"
	gormigrate "github.com/go-gormigrate/gormigrate/v2"
	logging "github.com/ipfs/go-log/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	log = logging.Logger("router")
	DB  *gorm.DB
)

func init() {
	initializers.LoadEnvVariables()
}

func ConnectToDB() {
	dbDsn := os.Getenv("DB_URL")

	var err error
	debug, err := strconv.ParseBool(os.Getenv("DRY_RUN"))

	if err != nil {
		log.Fatal("Failed to get or parse DRY_RUN env variable")
	}

	DB, err = OpenDatabase(dbDsn, debug)

	if err != nil {
		log.Fatal("Failed to connext to db")
	}
}

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

	if err != nil {
		return nil, err
	}

	if err = m.Migrate(); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}
	log.Debugf("Migration ran successfully")

	return DB, nil
}
