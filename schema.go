package main

import (
	"github.com/ipfs/go-cid"
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

	Policy Policy `json:"policy"`

	Collections      []Collection      `json:"collections"`
	ConstraintLabels []ConstraintLabel `json:"constraint_labels"`
	ValueLabels      []ValueLabel      `json:"value_labels"`
	StorageProviders []StorageProvider `json:"storage_providers" gorm:"many2many:tenant_storage_providers;"`

	MaximumBytesInFlight uint `json:"maximum_bytes_in_flight"` // Global maximum for all the tenant's content
	Suspended            bool `json:"suspended"`               // Tenant is suspended
}

type Policy struct {
	gorm.Model
	// Eligibility defines whether or not the SP is offered this Storage Contract by the broker
	Eligibility     []Clause
	AutoApprove     bool            `json:"auto_approve"` // If true, SPs can subscribe without approval
	StorageContract StorageContract `json:"storage_contract"`
}

// The Storage Contract CID is signed by the SP when they subscribe via the Deal Broker
type StorageContract struct {
	gorm.Model
	StorageContractCID cid.Cid `json:"storage_contract_cid"`
	Content            struct {
		InfoLines []string `json:"info_lines"`
	}
	Retrieval struct {
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
	gorm.Model
	BasePolicyID uint               `json:"base_policy_id"`
	Attribute    string             `json:"attribute"`
	Operator     ComparisonOperator `json:"operator"`
	Value        string             `json:"value"`
}

type Collection struct {
	gorm.Model
	ContentIndexUri string `json:"content_index_uri"`
	// ? An auth header (i.e, ContentIndexAuth) may be required if the index endpoint is protected
	ReplicationConstraints ReplicationConstraint `json:"replication_constraints"`
}

type ReplicationConstraint struct {
	gorm.Model
	CollectionID uint `json:"collection_id"`
	TenantID     uint `json:"tenant_id"`

	ConstraintID  uint `json:"constraint_id"`
	ConstraintMax uint `json:"constraint_max"`
}

type TenantStorageProviders struct {
	TenantID             uint `json:"tenant_id" gorm:"primaryKey"`
	SPID                 SPID `json:"spid" gorm:"primaryKey"`
	Suspended            bool `json:"suspended"`
	Approved             bool `json:"approved"`
	MaximumBytesInFlight uint `json:"maximum_bytes_in_flight"` // Maximum bytes this SP can have in flight from the tenant
}

type StorageProvider struct {
	gorm.Model
	SPID             SPID              `json:"spid" gorm:"primaryKey"`
	Tenants          []Tenant          `json:"tenants" gorm:"many2many:tenant_storage_providers;"`
	ConstraintValues []ConstraintValue `json:"constraint_values"` // Computed from ExternalValidationService
}

type ConstraintLabel struct {
	gorm.Model
	TenantID     uint   `json:"tenant_id"`
	ConstraintID uint   `json:"constraint_id"`
	Label        string `json:"label"`
}

// ! unique constraint, for each tenant/classID and tenant/label must be unique
// db.Table("class_labels").AddUniqueIndex("idx_class_labels_tenant_id_class_id", "tenant_id", "class_id")
// db.Table("class_labels").AddUniqueIndex("idx_class_labels_tenant_id_label", "tenant_id", "label")

type ValueLabel struct {
	gorm.Model
	TenantID uint   `json:"tenant_id"`
	ValueID  uint   `json:"value_id"`
	Label    string `json:"label"`
}

// ! unique constraint, for each tenant/ValueID and tenant/label must be unique
// db.Table("value_labels").AddUniqueIndex("idx_value_labels_tenant_id_class_id", "tenant_id", "class_id")
// db.Table("value_labels").AddUniqueIndex("idx_value_labels_tenant_id_label", "tenant_id", "label")

type ConstraintValue struct {
	ConstraintID uint `json:"constraint_id"`
	Value        uint `json:"value"`
}
