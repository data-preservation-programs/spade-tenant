package db

import (
	"database/sql/driver"
	"errors"
	"time"

	"github.com/jackc/pgtype"
	"gorm.io/gorm"
)

// TODO: apply not null to fields where reasonable

// todo: change to enum type if easy to do
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

type ID int32

type ModelBase struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Tenant struct {
	ModelBase
	TenantID ID `json:"tenant_id" gorm:"primaryKey"`

	Collections []Collection `json:"collections"`
	Labels      []Label      `json:"labels"`
	SPs         []SP         `json:"storage_providers" gorm:"many2many:tenants_sps;"`

	TenantAddresses []Address `json:"tenant_addresses"`

	TenantStorageContractCid string `json:"tenant_storage_contract_cid" gorm:"column:tenant_storage_contract_cid"`

	TenantSpEligibility []TenantSPEligibilityClauses

	// Meta:
	// - tenant_suspended
	TenantMeta pgtype.JSONB `gorm:"type:jsonb;default:'{}';not null"`

	// Settings:
	// - SP Auto Approve
	// - SP Auto Suspend
	// - Max In Flight GiB
	TenantSettings pgtype.JSONB `gorm:"type:jsonb;default:'{}';not null"`

	// add: policy_storage_contract_cid

	// sp_eligibility_criteria []clause
}

type Address struct {
	ModelBase
	TenantID         ID     `json:"tenant_id" gorm:"uniqueIndex:idx_address_tenant_id;not null"`
	Address          string `json:"address" gorm:"uniqueIndex:idx_address_tenant_id;not null"`
	AddressActorID   uint   `json:"actor_id"`
	AddressIsSigning bool   `json:"is_signing" gorm:"default:true;not null"`
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

// Generic SP <-> Tenant Eligibility Clause, specified as a `attribute`, `operator` and `value`
// Attribute is formatted as a path, i.e location.city, retrieval.success_rate
// examples:
// location.country ComparisonOperator.IncludedIn [CAN, USA]
// retrieval.success_rate ComparisonOperator.GreaterThan 0.98
type TenantSPEligibilityClauses struct {
	ModelBase
	TenantID        ID                 `json:"tenant_id"`
	ClauseAttribute string             `json:"attribute"`
	ClauseOperator  ComparisonOperator `json:"operator"`
	ClauseValue     string             `json:"value"`
}

type Collection struct {
	ModelBase
	CollectionID          ID           `json:"collection_id" gorm:"primaryKey"`
	TenantID              ID           `json:"tenant_id"`
	CollectionName        string       `json:"collection_name"`
	CollectionActive      bool         `json:"collection_active"`
	CollectionPieceSource pgtype.JSONB `gorm:"type:jsonb;default:'{}';not null"`
	CollectionDealParams  pgtype.JSONB `gorm:"type:jsonb;default:'{}';not null"`

	ReplicationConstraints []ReplicationConstraint `json:"replication_constraints"`
}

type ReplicationConstraint struct {
	ModelBase
	CollectionID ID `json:"collection_id"`

	ConstraintID  ID   `json:"constraint_id"`
	ConstraintMax uint `json:"constraint_max"`
}

type TenantSpState string

const (
	TenantSpStateEligible  TenantSpState = "eligible"
	TenantSpStatePending   TenantSpState = "pending"
	TenantSpStateActive    TenantSpState = "active"
	TenantSpStateSuspended TenantSpState = "suspended"
)

func (s *TenantSpState) Scan(value interface{}) error {
	strVal, ok := value.(string)
	if !ok {
		return errors.New("failed to scan TenantSpState")
	}
	*s = TenantSpState(strVal)
	return nil
}

func (s *TenantSpState) Value() (driver.Value, error) {
	return string(*s), nil
}

// Many:Many relation table between Tenants and SPs
// enum('eligible', 'pending', 'active', 'suspended')
type TenantsSPs struct {
	TenantID      ID            `json:"tenant_id" gorm:"primaryKey"`
	SPID          ID            `json:"sp_id" gorm:"primaryKey;column:'sp_id'"`
	TenantSpState TenantSpState `gorm:"type:tenant_sp_state;default:'eligible';not null"`
	// TODO: Enum type - test it out
	TenantSpsMeta pgtype.JSONB `gorm:"type:jsonb;default:'{}';not null"`
}

func (TenantsSPs) TableName() string {
	return "tenants_sps"
}

type SP struct {
	ModelBase
	SPID    ID       `json:"sp_id" gorm:"primaryKey"`
	Tenants []Tenant `json:"tenants" gorm:"many2many:tenants_sps;"`
}

// A label maps a uint to a human readable string
// It is used for both constraints (i.e, location.country) and values (i.e, CANADA)
// Each Tenant has their own unique set of labels
type Label struct {
	TenantID     ID           `json:"tenant_id" gorm:"uniqueIndex:idx_label_tenant_id_label_id;uniqueIndex:idx_label_tenant_id_label_text"`
	LabelID      ID           `json:"id" gorm:"uniqueIndex:idx_label_tenant_id_label_id"`
	LabelText    string       `json:"label" gorm:"uniqueIndex:idx_label_tenant_id_label_text"`
	LabelOptions pgtype.JSONB `gorm:"type:jsonb;default:'{}';not null"`
}
