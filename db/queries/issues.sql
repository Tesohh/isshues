-- name: GetIssueFromId :one
SELECT * FROM issues
WHERE id = $1;

-- name: GetIssueFromCode :one
SELECT * FROM issues
WHERE code = $1 
AND project_id = $2;

-- name: InsertIssue :one
INSERT INTO issues (title, code, description, status, priority, project_id, recruiter_user_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;
 
-- name: GetIssuesCountInProject :one
SELECT COUNT(*) FROM issues
WHERE project_id = $1;

-- name: BulkInsertIssueAssignees :copyfrom
INSERT INTO issue_assignees (issue_id, user_id) VALUES ($1, $2);

-- name: BulkInsertIssueLabels :copyfrom
INSERT INTO issue_labels (issue_id, label_id) VALUES ($1, $2);

-- name: BulkInsertIssueRelationships :copyfrom
INSERT INTO issue_relationships (from_issue_id, to_issue_id, category) VALUES ($1, $2, $3);
