
-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE customers (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    revision    INT NOT NULL DEFAULT 1,  -- tracks the latest revision number
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- other fields representing the latest revision
    name        VARCHAR(255) NOT NULL
);

CREATE TABLE customer_revisions (
    customer_id     UUID REFERENCES customers(id),
    revision        INT NOT NULL,
    created_at      TIMESTAMP NOT NULL,
    -- other fields representing the revision
    name            VARCHAR(255) NOT NULL,
    PRIMARY KEY (customer_id, revision)
);

-- +migrate Down

DROP TABLE customer_revisions;
DROP TABLE customers;
