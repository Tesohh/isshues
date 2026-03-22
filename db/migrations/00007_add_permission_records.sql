-- +goose Up
INSERT INTO global_permissions (id, description) VALUES
('create-projects', 'Create Projects');

-- Project permissions, eg. `write-issues`, `read-issues`, `edit-project`, `delete-project` are given to the group.
INSERT INTO project_permissions (id, description) VALUES
('write-issues', 'Create and edit issues'),
('read-issues', 'Read issues'),
('edit-project', 'Edit project'),
('delete-project', 'Delete project');


-- +goose Down
DELETE FROM global_permissions
WHERE id IN ('create-projects');

DELETE FROM project_permissions
WHERE id IN ('write-issues', 'read-issues', 'edit-project', 'delete-project');
