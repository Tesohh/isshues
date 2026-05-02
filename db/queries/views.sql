-- name: GetViewById :one
SELECT * FROM views
WHERE id = $1;

-- name: GetAllViewsInProject :many
SELECT * FROM views
WHERE project_id = $1;

-- name: InsertView :one
INSERT INTO views (project_id, name, title, 
	statuses, 
	priority, priority_mode, 
	labels_mode, 
	assignees_mode, assignees_include_viewer,
	assignee_groups_mode,
	sort_by, sort_order, 
	style)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
RETURNING id;

-- name: BulkInsertViewAssignees :copyfrom
INSERT INTO view_assignees (view_id, user_id) VALUES ($1, $2);

-- name: BulkInsertViewGroupAssignees :copyfrom
INSERT INTO view_group_assignees (view_id, group_id) VALUES ($1, $2);

-- name: BulkInsertViewLabels :copyfrom
INSERT INTO view_labels (view_id, label_id) VALUES ($1, $2);
