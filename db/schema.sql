CREATE TABLE IF NOT EXISTS "migrations" (
    "id" varchar(255),
    PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "tenants" (
    "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz NOT NULL,
    "deleted_at" timestamptz,
    "tenant_id" serial,
    "tenant_storage_contract_cid" text NOT NULL,
    "tenant_meta" jsonb NOT NULL DEFAULT '{}',
    "tenant_settings" jsonb NOT NULL DEFAULT '{}',
    PRIMARY KEY ("tenant_id")
);

CREATE INDEX IF NOT EXISTS "idx_tenants_deleted_at" ON "tenants" ("deleted_at");

CREATE TABLE IF NOT EXISTS "addresses" (
    "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz NOT NULL,
    "deleted_at" timestamptz,
    "tenant_id" integer,
    "address_robust" text,
    "address_actor_id" bigint,
    "address_is_signing" boolean NOT NULL DEFAULT true,
    PRIMARY KEY ("tenant_id", "address_robust"),
    CONSTRAINT "fk_tenants_tenant_addresses" FOREIGN KEY ("tenant_id") REFERENCES "tenants" ("tenant_id")
);

CREATE INDEX IF NOT EXISTS "idx_addresses_deleted_at" ON "addresses" ("deleted_at");

CREATE TABLE IF NOT EXISTS "tenant_sp_eligibility_clauses" (
    "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz NOT NULL,
    "deleted_at" timestamptz,
    "tenant_id" integer,
    "clause_attribute" text,
    "clause_operator" comparison_operator NOT NULL,
    "clause_value" jsonb NOT NULL,
    PRIMARY KEY ("tenant_id","clause_attribute"),
    CONSTRAINT "fk_tenants_tenant_sp_eligibility" FOREIGN KEY ("tenant_id") REFERENCES "tenants"("tenant_id")
);

CREATE INDEX IF NOT EXISTS "idx_tenant_sp_eligibility_clauses_deleted_at" ON "tenant_sp_eligibility_clauses" ("deleted_at");

CREATE TABLE IF NOT EXISTS "collections" (
    "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz NOT NULL,
    "deleted_at" timestamptz,
    "collection_id" uuid,
    "tenant_id" integer NOT NULL,
    "collection_name" text NOT NULL,
    "collection_active" boolean NOT NULL,
    "collection_piece_source" jsonb NOT NULL DEFAULT '{}',
    "collection_deal_params" jsonb NOT NULL DEFAULT '{}',
    PRIMARY KEY ("collection_id"),
    CONSTRAINT "fk_tenants_collections" FOREIGN KEY ("tenant_id") REFERENCES "tenants" ("tenant_id")
);

CREATE INDEX IF NOT EXISTS "idx_collections_deleted_at" ON "collections" ("deleted_at");

CREATE TABLE IF NOT EXISTS "labels" (
    "tenant_id" integer NOT NULL,
    "label_id" integer NOT NULL,
    "label_text" text NOT NULL,
    "label_options" jsonb NOT NULL DEFAULT '{}',
    CONSTRAINT "fk_tenants_labels" FOREIGN KEY ("tenant_id") REFERENCES "tenants" ("tenant_id")
);

CREATE UNIQUE INDEX IF NOT EXISTS "idx_tenant_id_label_text" ON "labels" ("tenant_id", "label_text");
CREATE UNIQUE INDEX IF NOT EXISTS "idx_tenant_id_label_id" ON "labels" ("tenant_id", "label_id");

CREATE TABLE IF NOT EXISTS "sps" (
    "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz NOT NULL,
    "deleted_at" timestamptz,
    "sp_id" serial,
    PRIMARY KEY ("sp_id")
);

CREATE INDEX IF NOT EXISTS "idx_sps_deleted_at" ON "sps" ("deleted_at");

CREATE TABLE IF NOT EXISTS "tenants_sps" (
    "tenant_id" integer,
    "sp_id" integer,
    "tenant_sp_state" tenant_sp_state NOT NULL DEFAULT 'eligible',
    "tenant_sp_state_info" text,
    "tenant_sps_meta" jsonb NOT NULL DEFAULT '{}',
    PRIMARY KEY ("tenant_id", "sp_id"),
    CONSTRAINT "fk_tenants_s_ps" FOREIGN KEY ("tenant_id") REFERENCES "tenants" ("tenant_id"),
    CONSTRAINT "fk_sps_tenants" FOREIGN KEY ("sp_id") REFERENCES "sps" ("sp_id")
);

CREATE TABLE IF NOT EXISTS "replication_constraints" (
    "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz NOT NULL,
    "deleted_at" timestamptz,
    "collection_id" uuid,
    "tenant_id" integer NOT NULL,
    "constraint_id" integer,
    "constraint_max" bigint NOT NULL,
    PRIMARY KEY ("collection_id", "constraint_id"),
    CONSTRAINT "fk_collections_replication_constraints" FOREIGN KEY ("collection_id") REFERENCES "collections" ("collection_id")
);

CREATE INDEX IF NOT EXISTS "idx_replication_constraints_deleted_at" ON "replication_constraints" ("deleted_at");

CREATE TABLE IF NOT EXISTS "sp_attributes" (
    "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz NOT NULL,
    "deleted_at" timestamptz,
    "attribute_id" integer,
    "sp_id" integer,
    "tenant_id" integer,
    "attribute_label_id" integer,
    "attribute_value_id" integer,
    PRIMARY KEY ("attribute_id", "attribute_label_id", "attribute_value_id"),
    CONSTRAINT "fk_tenants_sps_sp_attributes" FOREIGN KEY ("tenant_id", "sp_id") REFERENCES "tenants_sps" ("tenant_id", "sp_id")
);

CREATE INDEX IF NOT EXISTS "idx_sp_attributes_deleted_at" ON "sp_attributes" ("deleted_at");