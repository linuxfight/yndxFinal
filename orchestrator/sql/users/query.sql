-- name: GetById :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetByName :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: Create :exec
INSERT INTO users (
    id, username, password_hash
) VALUES (
             $1, $2, $3
         )
    RETURNING *;

-- name: CreateSchema :exec
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(26) NOT NULL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL
);