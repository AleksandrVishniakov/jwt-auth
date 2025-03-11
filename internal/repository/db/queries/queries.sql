-- name: CreateUser :one
INSERT INTO users (login, password_hash, role_id)
VALUES (
    $1, 
    $2, 
    (SELECT id FROM roles WHERE is_default = true LIMIT 1)
)
RETURNING id;

-- name: CreateSuperUser :one
INSERT INTO users (login, password_hash, role_id)
VALUES (
    $1, 
    $2, 
    (SELECT id FROM roles WHERE is_super = true LIMIT 1)
)
RETURNING id;

-- name: GetUserById :one
SELECT u.*, r.alias, r.permissions_mask
FROM users u JOIN roles r
ON u.role_id = r.id
WHERE u.id = $1;

-- name: GetUserByLogin :one
SELECT u.*, r.alias, r.permissions_mask
FROM users u JOIN roles r
ON u.role_id = r.id
WHERE u.login = $1;




-- name: UpdateRoleById :exec
UPDATE users
SET role_id = (SELECT roles.id FROM roles WHERE roles.alias = $2 LIMIT 1)
WHERE users.id = $1;

-- name: UpsertRole :one
INSERT INTO roles (alias, is_default, is_super, permissions_mask)
VALUES ($1, $2, $3, $4)
ON CONFLICT (alias) 
DO UPDATE SET 
    is_default = EXCLUDED.is_default,
    is_super = EXCLUDED.is_super,
    permissions_mask = EXCLUDED.permissions_mask
RETURNING id;

-- name: GetRoleByAlias :one
SELECT * FROM roles
WHERE roles.alias = $1;