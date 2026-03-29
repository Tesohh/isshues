-- name: GetIssueFromId :one
SELECT * FROM issues
WHERE id = $1;

-- name: GetIssueFromCode :one
SELECT * FROM issues
WHERE code = $1 
AND project_id = $2;
