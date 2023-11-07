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

type Tenant struct {
	gorm.Model

	Policy Policy

	Collections      []Collection      `json:"collections"`
	Labels           []Label           `json:"labels"`
	StorageProviders []StorageProvider `json:"storage_providers" gorm:"many2many:tenant_storage_providers;"`

	TenantSuspended bool `json:"tenant_suspended"`

	TenantWallets   []Address `json:"tenant_wallets"`
	TenantAddresses []Address `json:"tenant_addresses"`

	TenantMeta pgtype.JSONB `gorm:"type:jsonb;default:'[]';not null"`

	// Settings:
	// - SP Auto Approve
	// - SP Auto Suspend
	// - Max In Flight GiB
	TenantSettings pgtype.JSONB `gorm:"type:jsonb;default:'[]';not null"`
}

type Address struct {
	gorm.Model
	TenantID  uint   `json:"tenant_id" gorm:"uniqueIndex:idx_address_tenant_id;"`
	Address   string `json:"address" gorm:"uniqueIndex:idx_address_tenant_id;"`
	ActorID   uint   `json:"actor_id"`
	IsSigning bool   `json:"is_signing" gorm:"default:true"`
}

type Policy struct {
	gorm.Model
	TenantID                 uint `json:"tenant_id"`
	PolicyEligibility        []Clause
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
	TenantID              uint         `json:"tenant_id"`
	CollectionName        string       `json:"collection_name"`
	CollectionActive      bool         `json:"collection_active"`
	CollectionPieceSource pgtype.JSONB `gorm:"type:jsonb;default:'[]';not null"`
	CollectionDealParams  pgtype.JSONB `gorm:"type:jsonb;default:'[]';not null"`

	ReplicationConstraints []ReplicationConstraint `json:"replication_constraints"`
}

type ReplicationConstraint struct {
	gorm.Model
	CollectionID uint `json:"collection_id"`

	ConstraintID  uint `json:"constraint_id"`
	ConstraintMax uint `json:"constraint_max"`
}

type TenantStorageProvider struct {
	TenantID                  uint         `json:"tenant_id" gorm:"primaryKey"`
	SPID                      SPID         `json:"spid" gorm:"primaryKey"`
	Suspended                 bool         `json:"suspended"`
	Approved                  bool         `json:"approved"`
	TenantStorageProviderMeta pgtype.JSONB `gorm:"type:jsonb;default:'[]';not null"`
}

type StorageProvider struct {
	SPID    SPID     `json:"spid" gorm:"primaryKey"`
	Tenants []Tenant `json:"tenants" gorm:"many2many:tenant_storage_providers;"`
}

// A label maps a uint to a human readable string
// It is used for both constraints (i.e, location.country) and values (i.e, CANADA)
// Each Tenant has their own unique set of labels
type Label struct {
	UUID         string         `json:"uuid" gorm:"primaryKey"`
	TenantID     uint           `json:"tenant_id" gorm:"uniqueIndex:idx_label_tenant_id_id;uniqueIndex:idx_labels_tenant_id_label"`
	ID           uint           `json:"id" gorm:"uniqueIndex:idx_label_tenant_id_id"`
	Label        string         `json:"label" gorm:"uniqueIndex:idx_labels_tenant_id_label"`
	LabelOptions []LabelOptions `json:"label_options" gorm:"foreignKey:LabelUUID"`
}

type LabelOptions struct {
	LabelUUID string `json:"label_uuid" gorm:"primaryKey"`
	Option    string `json:"option"`
	Value     uint   `json:"value"`
}
