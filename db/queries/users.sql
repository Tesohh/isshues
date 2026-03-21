-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: IsUsernameTaken :one
SELECT 1 FROM users
WHERE username = $1;

-- name: InsertUser :one
INSERT INTO users (username) VALUES ($1)
RETURNING id;
 
-- name: UserHasProjectPermission :one
SELECT EXISTS (
    SELECT 1 FROM group_permissions gp
    JOIN user_group_memberships mb ON mb.group_id = gp.group_id
    JOIN groups g ON mb.group_id = g.id
    WHERE mb.user_id = $1 
    AND gp.project_permission_id = $2
    AND g.project_id = $3
);

-- name: UserHasGlobalPermission :one
SELECT EXISTS (
    SELECT 1 FROM global_user_permissions
    WHERE user_id = $1 AND global_permission_id = $2
);
