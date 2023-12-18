package core

import (
	"log"

	"github.com/data-preservation-programs/spade-tenant/config"
	"github.com/data-preservation-programs/spade-tenant/db"
	logging "github.com/ipfs/go-log/v2"
	"gorm.io/gorm"
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

	dbi, err := db.OpenDatabase(config.DB_URL, config.DEBUG, config.DRY_RUN)
	if err != nil {
		log.Fatalf("could not open db: %s", err)
	}

	return &SpdTenantSvc{DB: dbi, DryRunMode: config.DRY_RUN, Config: config}
}
