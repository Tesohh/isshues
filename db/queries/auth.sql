-- name: GetUserIdFromSSHFingerprint :one
SELECT user_id FROM user_ssh_keys
WHERE fingerprint = $1;

-- name: RegisterSSHFingerprintToUser :exec
INSERT INTO user_ssh_keys (user_id, fingerprint) VALUES ($1, $2);
