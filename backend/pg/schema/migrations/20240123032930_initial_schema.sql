-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS idx_user (
    id              INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      VARCHAR NOT NULL,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by      VARCHAR NOT NULL,
    user_id         VARCHAR NOT NULL UNIQUE,      
    email           VARCHAR NOT NULL UNIQUE,      
    auth            VARCHAR NOT NULL,             
    first_name      VARCHAR NOT NULL,             
    last_name       VARCHAR NOT NULL,             
    title           VARCHAR NOT NULL,             
    props           JSONB NOT NULL,               
    -- perms           VARCHAR[] NOT NULL            
);

CREATE TABLE IF NOT EXISTS idx_group (
    id              INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      VARCHAR NOT NULL,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by      VARCHAR NOT NULL,

    name            VARCHAR NOT NULL UNIQUE, 
    service_id      INT NOT NULL,
    display_name    VARCHAR NOT NULL,
    description     VARCHAR NOT NULL,

    CONSTRAINT fk FOREIGN KEY

);

CREATE TABLE IF NOT EXISTS idx_service (
    id              INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      VARCHAR NOT NULL,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by      VARCHAR NOT NULL,

    name            VARCHAR NOT NULL,
    owner_id        INT NOT NULL,
    display_name    VARCHAR NOT NULL,
    permissions     JSONB NOT NULL,

    CONSTRAINT fk_service_owner FOREIGN KEY(owner_id)
        REFERENCES idx_service(id) ON DELETE CASCADEs
);

CREATE TABLE IF NOT EXISTS user_pass(
    id INT PRIMARY KEY,
    password_hash VARCHAR NOT NULL,
    CONSTRAINT fk_user_password FOREIGN KEY(id)
            REFERENCES idx_user(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS idx_event(
	id              BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	op  		    VARCHAR NOT NULL,
    ev_type         VARCHAR NOT NULL,
	user_id		    VARCHAR NOT NULL,
	created_on		TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	errors          VARCHAR[]  DEFAULT '{}',
	metadata		JSONB
);

CREATE TABLE IF NOT EXISTS user_to_group (
    user_id         INT NOT NULL,
    group_id        INT NOT NULL,

    PRIMARY KEY(user_id, group_id),
    CONSTRAINT fk_g2u_group FOREIGN KEY(group_id) REFERENCES idx_group(id)
    CONSTRAINT fk_g2u_user FOREIGN KEY(user_id) REFERENCES idx_user(id),
        ON DELETE CASCADE
)

-- CREATE TABLE IF NOT EXISTS user_to_perm (
--     user_id         INT NOT NULL,
--     perm_id         VARCHAR NOT NULL,

--     PRIMARY KEY(user_id, perm_id),
--     CONSTRAINT fk_u2p_user FOREIGN KEY(user_id) REFERENCES idx_user(id)
--     ON DELETE CASCADE
-- )

CREATE TABLE IF NOT EXISTS group_to_perm (
    group_id         INT NOT NULL,
    perm_id        VARCHAR NOT NULL,

    PRIMARY KEY(group_id, perm_id),
    CONSTRAINT fk_g2p_group FOREIGN KEY(group_id) REFERENCES idx_group(id)
        ON DELETE CASCADE
)

-- CREATE TABLE IF NOT EXISTS service_to_group (
--     service_id         INT NOT NULL,
--     group_id        VARCHAR NOT NULL,

--     PRIMARY KEY(service_id, group_id),
--     CONSTRAINT fk_s2g_service FOREIGN KEY(service_id) REFERENCES idx_service(id)
--         ON DELETE CASCADE,
--     CONSTRAINT fk_s2g_group FOREIGN KEY(group_id) REFERENCES idx_group(id)
--         ON DELETE CASCADE
-- )

CREATE TABLE IF NOT EXISTS credential (
    id VARCHAR,
    item_type VARCHAR
    password_hash VARCHAR NOT NULL,
    PRIMARY KEY(id, item_type)
)

-- CREATE TABLE IF NOT EXISTS service_secret (
--     id INT PRIMARY KEY,
--     password_hash VARCHAR NOT NULL,
--     CONSTRAINT fk_s2p FOREIGN KEY(id) REFERENCES idx_service(id)
--         ON DELETE CASCADE
-- )

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- DROP TABLE user_to_perm
-- DROP TABLE service_to_pw
DROP TABLE credential
DROP TABLE service_to_group
DROP TABLE group_to_perm
DROP TABLE user_to_group
DROP TABLE idx_event
DROP TABLE user_pass
DROP TABLE idx_service
DROP TABLE idx_group
DROP TABLE idx_user
-- +goose StatementEnd
