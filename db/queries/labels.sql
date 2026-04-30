-- name: GetLabelFromId :one
SELECT * FROM labels
WHERE id = $1;

-- name: GetLabelFromName :one
SELECT * FROM labels
WHERE name = $1
AND project_id = $2;

-- name: GetLabelsByIDBulk :many
SELECT * FROM labels
WHERE id = ANY($1::bigint[]);

-- name: InsertLabelBasic :one
INSERT INTO labels (name, project_id)
VALUES ($1, $2)
RETURNING *;

-- name: InsertLabel :one
INSERT INTO labels (name, symbol, color_key, project_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

