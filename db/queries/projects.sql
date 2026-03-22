-- name: GetUserProjectMemberships :many
SELECT DISTINCT p.* FROM projects p
JOIN groups g ON g.project_id = p.id
JOIN user_group_memberships mb ON mb.group_id = g.id
WHERE mb.user_id = $1;
