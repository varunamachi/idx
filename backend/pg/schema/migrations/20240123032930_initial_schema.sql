-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user (
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
    props           JSONB,                       
    perms           VARCHAR[]            
);

CREATE TABLE IF NOT EXISTS pgroup (
    id              INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      VARCHAR NOT NULL,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by      VARCHAR NOT NULL,

    name            VARCHAR NOT NULL UNIQUE, 
    display_name    VARCHAR NOT NULL,
    description     VARCHAR NOT NULL,
    perms           VARCHAR[]
);

CREATE TABLE IF NOT EXISTS app (

);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
