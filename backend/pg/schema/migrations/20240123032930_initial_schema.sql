-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS idx_user (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    created_on TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by VARCHAR NOT NULL,
    updated_on TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by VARCHAR NOT NULL,
    user_name VARCHAR NOT NULL UNIQUE,
    email VARCHAR NOT NULL UNIQUE,
    auth VARCHAR NOT NULL,
    state VARCHAR NOT NULL,
    first_name VARCHAR NOT NULL,
    last_name VARCHAR NOT NULL,
    title VARCHAR NOT NULL,
    props JSONB NOT NULL
);

CREATE TABLE IF NOT EXISTS idx_service (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    created_on TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by VARCHAR NOT NULL,
    updated_on TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by VARCHAR NOT NULL,
    name VARCHAR NOT NULL UNIQUE,
    owner_id INT NOT NULL,
    display_name VARCHAR NOT NULL,
    permissions JSONB NOT NULL,
    CONSTRAINT fk_service_owner FOREIGN KEY(owner_id) 
        REFERENCES idx_user(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS idx_group (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    created_on TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by VARCHAR NOT NULL,
    updated_on TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by VARCHAR NOT NULL,
    name VARCHAR NOT NULL UNIQUE,
    service_id INT NOT NULL,
    display_name VARCHAR NOT NULL,
    description VARCHAR NOT NULL,
    CONSTRAINT fk_service FOREIGN KEY(service_id) 
        REFERENCES idx_service(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS service_to_owner (
    service_id INT NOT NULL,
    admin_id INT NOT NULL,
    PRIMARY KEY(service_id, admin_id),
    CONSTRAINT fk_service_id FOREIGN KEY(service_id) 
        REFERENCES idx_service(id) ON DELETE CASCADE,
    CONSTRAINT fk_admin_admin FOREIGN KEY(admin_id) 
        REFERENCES idx_user(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS user_pass(
    id INT PRIMARY KEY,
    password_hash VARCHAR NOT NULL,
    CONSTRAINT fk_user_password FOREIGN KEY(id) 
        REFERENCES idx_user(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS idx_event(
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    op VARCHAR NOT NULL,
    ev_type VARCHAR NOT NULL,
    user_id INT NOT NULL,
    created_on TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    errors VARCHAR [] DEFAULT '{}',
    metadata JSONB
);

CREATE TABLE IF NOT EXISTS user_to_group (
    user_id INT NOT NULL,
    group_id INT NOT NULL,
    PRIMARY KEY(user_id, group_id),
    CONSTRAINT fk_g2u_group FOREIGN KEY(group_id) 
        REFERENCES idx_group(id) ON DELETE CASCADE,
    CONSTRAINT fk_g2u_user FOREIGN KEY(user_id) 
        REFERENCES idx_user(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS group_to_perm (
    group_id INT NOT NULL,
    perm_id VARCHAR NOT NULL,
    PRIMARY KEY(group_id, perm_id),
    CONSTRAINT fk_g2p_group FOREIGN KEY(group_id) 
        REFERENCES idx_group(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS credential (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    unique_name VARCHAR,
    item_type VARCHAR,
    password_hash VARCHAR NOT NULL,
    created_on TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    num_failed_auth INTEGER NOT NULL DEFAULT 0,

    -- last_failed_auth as default as now() to have a sane defalt
    -- since this field is used only when num_failed_auth is not 0
    -- and when the first failure sets this field, it should be fine
    last_failed_on TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    prev_passwords VARCHAR [] DEFAULT '{}',
    PRIMARY KEY(unique_name, item_type)
);


CREATE TABLE IF NOT EXISTS credential_policy (
    item_type    VARCHAR NOT NULL,
    pattern VARCHAR NOT NULL,
    expiry  TIME NOT NULL,
    max_retries INTEGER NOT NULL,
    retry_reset_days INTEGER NOT NULL DEFAULT 1,
    max_reuse INTEGER NOT NULL,
    PRIMARY KEY(item_type)
)

CREATE TABLE IF NOT EXISTS idx_token (
    token VARCHAR NOT NULL,
    unique_name VARCHAR NOT NULL,
    assoc_type VARCHAR NOT NULL,
    -- user or service
    operation VARCHAR NOT NULL,
    created_on TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY(token),
    UNIQUE(token, unique_name, assoc_type, operation)
);

-- TODO - Add credential initial policy

-- +goose StatementEnd
--
-- +goose Down
-- +goose StatementBegin
DROP TABLE idx_token;

DROP TABLE credential_policy;

DROP TABLE credential;

DROP TABLE group_to_perm;

DROP TABLE user_to_group;

DROP TABLE idx_event;

DROP TABLE user_pass;

DROP TABLE service_to_owner;

DROP TABLE idx_group;

DROP TABLE idx_service;

DROP TABLE idx_user;

--
-- +goose StatementEnd