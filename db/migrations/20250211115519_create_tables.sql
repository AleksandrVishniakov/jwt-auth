-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    alias VARCHAR(64) UNIQUE NOT NULL,
    is_default BOOLEAN NOT NULL DEFAULT false,
    is_super BOOLEAN NOT NULL DEFAULT false,
    permissions_mask BIGINT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    login VARCHAR(128) UNIQUE NOT NULL,
    role_id SERIAL REFERENCES roles(id) NOT NULL,
    password_hash VARCHAR(256) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CHECK ( length(login) >= 3 )
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
