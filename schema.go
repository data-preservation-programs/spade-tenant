package main

import (
	"time"

	"github.com/ipfs/go-cid"
	"github.com/jackc/pgtype"
	"gorm.io/gorm"
)

type SPID uint

type ComparisonOperator string

const (
	GreaterThan        ComparisonOperator = ">"
	LessThan           ComparisonOperator = "<"
	EqualTo            ComparisonOperator = "="
	GreaterThanOrEqual ComparisonOperator = ">="
	LessThanOrEqual    ComparisonOperator = "<="
	IncludedIn         ComparisonOperator = "in"
	ExcludedFrom       ComparisonOperator = "nin"
	NotEqualTo         ComparisonOperator = "!="
)

type OrmModel struct {
	ID        int32 `json:"id" gorm:"primaryKey"` // overwrite uint -> int32
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// gorm.Model
// may need to change some of the defaults (i.e, ID) to allow for more performance (constraints with the broker)

// Gorm -> prefix column name with table name for ease of querying later (where reasonable)
type Tenant struct {
	OrmModel

	Policy Policy

	Collections      []Collection      `json:"collections"`
	Labels           []Label           `json:"labels"`
	StorageProviders []StorageProvider `json:"storage_providers" gorm:"many2many:tenant_storage_providers;"`

	TenantDefaultMaxBytesInFlight uint
	TenantSuspended               bool `json:"tenant_suspended"` // Tenant is suspended
}

type Policy struct {
	OrmModel
	TenantID              int32 `json:"tenant_id"`
	PolicyEligibility     []Clause
	PolicyAutoApprove     bool            `json:"auto_approve"` // If true, SPs can subscribe without approval
	PolicyStorageContract StorageContract `json:"storage_contract" gorm:"embedded"`
}

// The Storage Contract CID is signed by the SP when they subscribe via the Deal Broker
type StorageContract struct {
	StorageContractCID cid.Cid `json:"storage_contract_cid"`

	// ? Do we need to store it in the DB? Or just force it to be in IPFS and use the CID only.
	StorageContractContent struct {
		InfoLines []string `json:"info_lines"`
	}
	StorageContentRetrieval struct {
		Mechanisms struct {
			IpldBitswap  bool `json:"ipld_bitswap"`
			Piece_Rrhttp bool `json:"piece_rrhttp"`
		}
		Sla struct {
			InfoLines []string `json:"info_lines"`
		}
	}
}

// A generic element of a policy, specified as a `attribute`, `operator` and `value`
// Attribute is formatted as a path, i.e location.city, retrieval.success_rate
// Some examples:
// location.country ComparisonOperator.IncludedIn [CAN, USA]
// retrieval.success_rate ComparisonOperator.GreaterThan 0.98
type Clause struct {
	OrmModel
	PolicyID        uint               `json:"policy_id"`
	ClauseAttribute string             `json:"attribute"`
	ClauseOperator  ComparisonOperator `json:"operator"`
	ClauseValue     string             `json:"value"`
}

type Collection struct {
	OrmModel
	TenantID               int32                   `json:"tenant_id"`
	CollectionPieceSource  pgtype.JSONB            `gorm:"type:jsonb;default:'[]';not null"`
	ReplicationConstraints []ReplicationConstraint `json:"replication_constraints"`
}

type ReplicationConstraint struct {
	OrmModel
	CollectionID uint `json:"collection_id"`

	ConstraintID  uint `json:"constraint_id"`
	ConstraintMax uint `json:"constraint_max"`

	ConstraintValues []ConstraintValue `json:"constraint_values" gorm:"foreignKey:ConstraintID"`
}

type TenantStorageProviders struct {
	TenantID             int32 `json:"tenant_id" gorm:"primaryKey"`
	SPID                 SPID  `json:"spid" gorm:"primaryKey"`
	Suspended            bool  `json:"suspended"`
	Approved             bool  `json:"approved"`
	MaximumBytesInFlight uint  `json:"maximum_bytes_in_flight"` // Maximum bytes this SP can have in flight from the tenant
}

type StorageProvider struct {
	SPID             SPID              `json:"spid" gorm:"primaryKey"`
	Tenants          []Tenant          `json:"tenants" gorm:"many2many:tenant_storage_providers;"`
	ConstraintValues []ConstraintValue `json:"constraint_values" gorm:"foreignKey:SPID"` // Computed from ExternalValidationService
}

// A label maps a uint to a human readable string
// It is used for both constraints (i.e, location.country) and values (i.e, CANADA)
// Each Tenant has their own unique set of labels
type Label struct {
	TenantID uint   `json:"tenant_id"`
	ID       uint   `json:"id"`
	Label    string `json:"label"`
}

// ! unique constraint, for each tenant/ValueID and tenant/label must be unique
// db.Table("labels").AddUniqueIndex("idx_labels_tenant_id_id", "tenant_id", "id")
// db.Table("labels").AddUniqueIndex("idx_labels_tenant_id_label", "tenant_id", "label")

type ConstraintValue struct {
	SPID         SPID `json:"spid" gorm:"primaryKey"`
	ConstraintID uint `json:"constraint_id"`
	Value        uint `json:"value"`
}
