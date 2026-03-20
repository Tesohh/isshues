-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: IsUsernameTaken :one
SELECT 1 FROM users
WHERE username = $1;

-- name: InsertUser :one
INSERT INTO users (username, is_admin) VALUES ($1, $2)
RETURNING id;
 
