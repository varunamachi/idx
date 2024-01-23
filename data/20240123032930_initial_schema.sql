-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user (

);

CREATE TABLE IF NOT EXISTS pgroup (

);

CREATE TABLE IF NOT EXISTS app (

);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
