-- name: GetUserIdFromAuth :one
SELECT user_ssh_keys.user_id FROM user_ssh_keys
JOIN users ON users.id = user_ssh_keys.user_id
WHERE fingerprint = $1 AND username = $2;

-- name: RegisterSSHFingerprintToUser :exec
INSERT INTO user_ssh_keys (user_id, fingerprint) VALUES ($1, $2);
