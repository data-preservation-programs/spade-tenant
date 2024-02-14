package db

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"gorm.io/gorm"
)

type ID int32

type ModelBase struct {
	CreatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Tenant struct {
	ModelBase
	TenantID ID `json:"tenant_id" gorm:"primaryKey"`

	Collections []Collection `json:"collections"`
	Labels      []Label      `json:"labels"`
	SPs         []TenantsSPs `json:"storage_providers"`

	TenantAddresses []Address `json:"tenant_addresses"`

	TenantStorageContractCid string `json:"tenant_storage_contract_cid" gorm:"column:tenant_storage_contract_cid;not null"`

	TenantSpEligibility []TenantSPEligibilityClauses

	// Meta:
	// - tenant_suspended
	TenantMeta pgtype.JSONB `gorm:"type:jsonb;default:'{}';not null"`

	TenantSettings pgtype.JSONB `gorm:"type:jsonb;default:'{}';not null"`
}

// JSON type for the TenantSettings field
type TenantSettings struct {
	SpAutoApprove  bool `json:"sp_auto_approve"`
	SpAutoSuspend  bool `json:"sp_auto_suspend"`
	MaxInFlightGiB uint `json:"max_in_flight_gib"`
}

type Address struct {
	ModelBase
	TenantID         ID      `json:"tenant_id" gorm:"primaryKey"`
	AddressRobust    *string `json:"address_robust" gorm:"primaryKey"`
	AddressActorID   *uint   `json:"actor_id"`
	AddressIsSigning *bool   `json:"is_signing" gorm:"default:true;not null"`
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

func (s *ComparisonOperator) Scan(value interface{}) error {
	strVal, ok := value.(string)
	if !ok {
		return errors.New("failed to scan TenantSpState")
	}
	*s = ComparisonOperator(strVal)
	return nil
}

func (s *ComparisonOperator) Value() (driver.Value, error) {
	return string(*s), nil
}

// Generic SP <-> Tenant Eligibility Clause, specified as a `attribute`, `operator` and `value`
// examples:
// location.country ComparisonOperator.IncludedIn [CAN, USA]
// retrieval.success_rate ComparisonOperator.GreaterThan 0.98
type TenantSPEligibilityClauses struct {
	ModelBase
	TenantID        ID                 `json:"tenant_id" gorm:"primaryKey"`
	ClauseAttribute string             `json:"attribute" gorm:"primaryKey"`
	ClauseOperator  ComparisonOperator `json:"operator" gorm:"type:comparison_operator;not null"`
	ClauseValue     string             `json:"value" gorm:"not null"`
}
type Collection struct {
	ModelBase
	CollectionID          uuid.UUID    `json:"collection_id" gorm:"type:uuid;primaryKey"`
	TenantID              ID           `json:"tenant_id" gorm:"not null"`
	CollectionName        *string      `json:"collection_name" gorm:"not null"`
	CollectionActive      *bool        `json:"collection_active" gorm:"not null"`
	CollectionPieceSource pgtype.JSONB `json:"collection_piece_source" gorm:"type:jsonb;default:'{}';not null"`
	CollectionDealParams  pgtype.JSONB `json:"collection_deal_params" gorm:"type:jsonb;default:'{}';not null"`

	ReplicationConstraints []ReplicationConstraint `json:"replication_constraints"`
}

type ReplicationConstraint struct {
	ModelBase
	CollectionID  uuid.UUID `json:"collection_id" gorm:"type:uuid;primaryKey"`
	TenantID      ID        `json:"tenant_id" gorm:"not null"`
	ConstraintID  ID        `json:"constraint_id" gorm:"primaryKey"`
	ConstraintMax int       `json:"constraint_max" gorm:"not null"`
}

type TenantSpState string

// TenantSpState represents the state of a Tenant's relationship with a Storage Provider
// The state machine is as follows:
// - eligible: A given storage provider is eligible (i.e, conforms to all Eligibility clauses) to be associated with a tenant
// - pending: A given storage provider has subscribed to the tenant's Storage Contract, and is awaiting approval from the tenant
// - active: A given storage provider has been approved by the tenant. This is the only state where deals may be made.
// - suspended: A given storage provider has been suspended by the tenant for breach of Storage Contract or related SLA terms. It may be moved back to Active state once it is rectified.
// - disabled: A given storage provider has been disabled. This is a terminal state.
const (
	TenantSpStateEligible  TenantSpState = "eligible"
	TenantSpStatePending   TenantSpState = "pending"
	TenantSpStateActive    TenantSpState = "active"
	TenantSpStateSuspended TenantSpState = "suspended"
	TenantSpStateDisabled  TenantSpState = "disabled"
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
type TenantsSPs struct {
	TenantID          ID            `json:"tenant_id" gorm:"primaryKey"`
	SPID              ID            `json:"sp_id" gorm:"primaryKey;column:sp_id"`
	TenantSpState     TenantSpState `gorm:"type:tenant_sp_state;column:tenant_sp_state;default:eligible;not null"`
	TenantSpStateInfo string        `json:"tenant_sp_state_info"`
	TenantSpsMeta     pgtype.JSONB  `gorm:"type:jsonb;default:'{}';not null"`
	SPAttributes      []SPAttribute `json:"sp_attributes" gorm:"foreignKey:TenantID,SPID;references:TenantID,SPID"`
}

// JSON type for the TenantSpsMeta field
type TenantsSpsMeta struct {
	Signature         string `json:"signature"`
	SignedContractCID string `json:"signed_contract_cid"`
}

func (tenantSPs *TenantsSPs) BeforeUpdate(tx *gorm.DB) error {

	var currentValue TenantsSPs
	tx.Model(&TenantsSPs{SPID: ID(tenantSPs.SPID), TenantID: tenantSPs.TenantID}).Find(&currentValue)

	// Valid state transitions:
	// eligilble -> pending
	if currentValue.TenantSpState == TenantSpStateEligible && tenantSPs.TenantSpState == TenantSpStatePending {
		return nil
	}
	// eligible -> active (auto-approve)
	if currentValue.TenantSpState == TenantSpStateEligible && tenantSPs.TenantSpState == TenantSpStateActive {
		return nil
	}
	// pending -> active
	if currentValue.TenantSpState == TenantSpStatePending && tenantSPs.TenantSpState == TenantSpStateActive {
		return nil
	}
	// active -> suspended
	if currentValue.TenantSpState == TenantSpStateActive && tenantSPs.TenantSpState == TenantSpStateSuspended {
		return nil
	}
	// suspended -> active
	if currentValue.TenantSpState == TenantSpStateSuspended && tenantSPs.TenantSpState == TenantSpStateActive {
		return nil
	}
	// *(any) -> disabled
	if tenantSPs.TenantSpState == TenantSpStateDisabled {
		return nil
	}

	return fmt.Errorf("cannot go from state %s to state %s", currentValue.TenantSpState, tenantSPs.TenantSpState)
}

func (TenantsSPs) TableName() string {
	return "tenants_sps"
}

type SP struct {
	ModelBase
	SPID    ID           `json:"sp_id" gorm:"primaryKey"`
	Tenants []TenantsSPs `json:"tenants"`
}

// A relation table between SPs and their attributes
// Each SP has a set of attributes, which are a combination of labels and values
// For example, an SP may have the following attributes:
// - location.country = CANADA
// - location.city = TORONTO
// This relationship is modelled as 2 foreign keys to the Labels table - one for the "attribute label" and one for the "value label"
type SPAttribute struct {
	ModelBase
	SPID             ID `json:"sp_id"`
	TenantID         ID `json:"tenant_id"`
	AttributeLabelID ID `json:"attribute_label_id" gorm:"primaryKey"`
	AttributeValueID ID `json:"attribute_value_id" gorm:"primaryKey"`
}

// A label maps a uint to a human readable string
// It is used for both constraints (i.e, location.country) and values (i.e, CANADA)
// Each Tenant has their own unique set of labels
type Label struct {
	TenantID     ID           `json:"tenant_id" gorm:"uniqueIndex:idx_tenant_id_label_id;uniqueIndex:idx_tenant_id_label_text;not null"`
	LabelID      ID           `json:"id" gorm:"uniqueIndex:idx_tenant_id_label_id;not null"`
	LabelText    string       `json:"label" gorm:"uniqueIndex:idx_tenant_id_label_text;not null"`
	LabelOptions pgtype.JSONB `gorm:"type:jsonb;default:'{}';not null"`
}
