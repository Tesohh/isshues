-- name: GetLabelFromId :one
SELECT * FROM labels
WHERE id = $1;

-- name: GetLabelFromName :one
SELECT * FROM labels
WHERE name = $1
AND project_id = $2;
