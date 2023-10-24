package db

import (
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

// type OrmModel struct {
// 	ID        int32 `json:"id" gorm:"primaryKey"` // overwrite uint -> int32
// 	CreatedAt time.Time
// 	UpdatedAt time.Time
// 	DeletedAt gorm.DeletedAt `gorm:"index"`
// }

type Tenant struct {
	gorm.Model

	Policy Policy

	Collections      []Collection      `json:"collections"`
	Labels           []Label           `json:"labels"`
	StorageProviders []StorageProvider `json:"storage_providers" gorm:"many2many:tenant_storage_providers;"`

	TenantDefaultMaxBytesInFlight uint
	TenantSuspended               bool `json:"tenant_suspended"` // Tenant is suspended

	TenantWallets             []Address `json:"tenant_wallets"`
	TenantAssociatedAddresses []Address `json:"tenant_associated_addresses"`

	TenantMeta pgtype.JSONB `gorm:"type:jsonb;default:'[]';not null"`
}

type Address struct {
	gorm.Model
	TenantID  uint   `json:"tenant_id" gorm:"uniqueIndex:idx_address_tenant_id;"`
	Address   string `json:"address" gorm:"uniqueIndex:idx_address_tenant_id;"`
	IsSigning bool   `json:"is_signing" gorm:"default:true"`
}
type Policy struct {
	gorm.Model
	TenantID                 uint `json:"tenant_id"`
	PolicyEligibility        []Clause
	PolicyAutoApprove        bool   `json:"auto_approve"` // If true, SPs can subscribe without approval
	PolicyStorageContractCID string `json:"storage_contract_cid"`
}

// * Schema of the Storage Contract which will be stored in IPFS, referenced by CID
// type StorageContract struct {
// StorageContractContent struct {
// 	InfoLines []string `json:"info_lines"`
// }
// StorageContentRetrieval struct {
// 	Mechanisms struct {
// 		IpldBitswap  bool `json:"ipld_bitswap"`
// 		Piece_Rrhttp bool `json:"piece_rrhttp"`
// 	}
// 	Sla struct {
// 		InfoLines []string `json:"info_lines"`
// 	}
// }
// }

// A generic element of a policy, specified as a `attribute`, `operator` and `value`
// Attribute is formatted as a path, i.e location.city, retrieval.success_rate
// Some examples:
// location.country ComparisonOperator.IncludedIn [CAN, USA]
// retrieval.success_rate ComparisonOperator.GreaterThan 0.98
type Clause struct {
	gorm.Model
	PolicyID        uint               `json:"policy_id"`
	ClauseAttribute string             `json:"attribute"`
	ClauseOperator  ComparisonOperator `json:"operator"`
	ClauseValue     string             `json:"value"`
}

type Collection struct {
	gorm.Model
	TenantID               int32                   `json:"tenant_id"`
	CollectionPieceSource  pgtype.JSONB            `gorm:"type:jsonb;default:'[]';not null"`
	ReplicationConstraints []ReplicationConstraint `json:"replication_constraints"`
}

type ReplicationConstraint struct {
	gorm.Model
	CollectionID uint `json:"collection_id"`

	ConstraintID  uint `json:"constraint_id"`
	ConstraintMax uint `json:"constraint_max"`

	ConstraintAttributes []ConstraintAttribute `json:"constraint_attributes" gorm:"foreignKey:ConstraintID"`
}

type TenantStorageProvider struct {
	TenantID                  int32        `json:"tenant_id" gorm:"primaryKey"`
	SPID                      SPID         `json:"spid" gorm:"primaryKey"`
	Suspended                 bool         `json:"suspended"`
	Approved                  bool         `json:"approved"`
	TenantStorageProviderMeta pgtype.JSONB `gorm:"type:jsonb;default:'[]';not null"`
	MaxBytesInFlight          uint         `json:"max_bytes_in_flight"` // Maximum bytes this SP can have in flight from the tenant
}

type StorageProvider struct {
	SPID                SPID                  `json:"spid" gorm:"primaryKey"`
	Tenants             []Tenant              `json:"tenants" gorm:"many2many:tenant_storage_providers;"`
	ConstraintAttribute []ConstraintAttribute `json:"constraint_attributes" gorm:"foreignKey:SPID"` // Computed from ExternalValidationService
}

// A label maps a uint to a human readable string
// It is used for both constraints (i.e, location.country) and values (i.e, CANADA)
// Each Tenant has their own unique set of labels
type Label struct {
	TenantID uint   `json:"tenant_id" gorm:"uniqueIndex:idx_label_tenant_id_id;uniqueIndex:idx_labels_tenant_id_label"`
	ID       uint   `json:"id" gorm:"uniqueIndex:idx_label_tenant_id_id"`
	Label    string `json:"label" gorm:"uniqueIndex:idx_labels_tenant_id_label"`
}

type ConstraintAttribute struct {
	SPID         SPID `json:"spid" gorm:"primaryKey"`
	ConstraintID uint `json:"constraint_id"`
	Value        uint `json:"value"`
}
