-- name: GetAll :many
SELECT * FROM expressions
ORDER BY id DESC;

-- name: GetById :one
SELECT * FROM expressions
WHERE id = $1 LIMIT 1;

-- name: GetByExpr :one
SELECT * FROM expressions
WHERE expr = $1 LIMIT 1;

-- name: Create :exec
INSERT INTO expressions (
    id, expr, res, finished, error
) VALUES (
             $1, $2, $3, $4, $5
         )
RETURNING *;

-- name: Update :exec
UPDATE expressions
SET res = $1,
    finished = $2,
    error = $3
WHERE id = $4;

-- name: CreateSchema :exec
CREATE TABLE IF NOT EXISTS expressions (
    id VARCHAR(26) NOT NULL PRIMARY KEY,
    expr TEXT NOT NULL UNIQUE,
    res FLOAT NOT NULL,
    finished BOOLEAN NOT NULL,
    error BOOLEAN NOT NULL
);