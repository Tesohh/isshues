-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUsersByIDBulk :many
SELECT * FROM users
WHERE id = ANY($1::bigint[]);

-- name: GetUserByUsernameLenient :one
SELECT * FROM users
WHERE shortname = $1
OR username LIKE '%$1%';

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1;

-- name: IsUsernameTaken :one
SELECT EXISTS (
    SELECT 1 FROM users
    WHERE username = $1
);

-- name: InsertUser :one
INSERT INTO users (username) VALUES ($1)
RETURNING id;

-- name: InsertUserSettings :one
INSERT INTO user_settings (user_id, theme) VALUES ($1, $2)
RETURNING user_id;
 
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

-- name: UserIsMemberOfProject :one
-- You are a member if you are in at least one group inside of a project.
SELECT EXISTS (
    SELECT 1 FROM user_group_memberships mb
    JOIN groups g ON mb.group_id = g.id
    WHERE mb.user_id = $1 
    AND g.project_id = $2
);

-- name: GetUserSettings :one
SELECT * FROM user_settings
WHERE user_id = $1;
