-- name: GetUserProjectMemberships :many
SELECT DISTINCT p.* FROM projects p
JOIN groups g ON g.project_id = p.id
JOIN user_group_memberships mb ON mb.group_id = g.id
WHERE mb.user_id = $1;

-- name: GetProjectByPrefix :one
SELECT * FROM projects
WHERE prefix = $1;

-- name: GetProjectById :one
SELECT * FROM projects
WHERE id = $1;

-- name: InsertProject :one
INSERT INTO projects (title, prefix)
VALUES ( $1, $2 )
RETURNING id;
