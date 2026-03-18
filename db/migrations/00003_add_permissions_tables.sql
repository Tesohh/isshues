-- +goose Up
CREATE TABLE global_permissions(
    id VARCHAR(24) PRIMARY KEY,
    description TEXT
);

CREATE TABLE global_user_permissions(
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    global_permission_id VARCHAR(24) REFERENCES global_permissions(id),
    PRIMARY KEY (user_id, global_permission_id)
);

CREATE TABLE project_permissions(
    id VARCHAR(24) PRIMARY KEY,
    description TEXT
);

CREATE TABLE project_user_permissions(
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    project_permission_id VARCHAR(24) REFERENCES project_permissions(id),
    PRIMARY KEY (user_id, project_permission_id)
);

-- +goose Down
DROP TABLE global_user_permissions;
DROP TABLE global_permissions;
DROP TABLE project_user_permissions;
DROP TABLE project_permissions;
