-- name: InsertGroup :one
INSERT INTO groups (name, color_key, mentionable, project_id)
VALUES ($1, $2, $3, $4)
RETURNING groups.id;

-- name: GrantPermissionToGroup :exec
INSERT INTO group_permissions (group_id, project_permission_id)
VALUES ($1, $2);

-- name: RemovePermissionFromGroup :exec
DELETE FROM group_permissions
WHERE group_id = $1 AND project_permission_id = $2;

-- name: AddUserToGroup :exec
INSERT INTO user_group_memberships (group_id, user_id)
VALUES ($1, $2);

-- name: RemoveUserFromGroup :exec
DELETE FROM user_group_memberships
WHERE group_id = $1 AND user_id = $2;
