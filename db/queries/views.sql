-- name: GetViewById :one
SELECT * FROM views
WHERE id = $1;

-- name: GetAllViewsInProject :many
SELECT * FROM views
WHERE project_id = $1;
