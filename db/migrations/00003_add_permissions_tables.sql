-- +goose Up
CREATE TABLE global_permissions(
    id VARCHAR(24) PRIMARY KEY,
    description TEXT
);

CREATE TABLE project_permissions(
    id VARCHAR(24) PRIMARY KEY,
    description TEXT
);


CREATE TABLE groups (
    id BIGSERIAL PRIMARY KEY,
    name TEXT, -- when null, consider this group as "anonymous"
    color CHAR(7), 
    mentionable BOOL NOT NULL,
    project_id BIGINT NOT NULL REFERENCES projects(id) ON DELETE CASCADE
);

CREATE TABLE group_permissions (
    group_id BIGINT NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    project_permission_id VARCHAR(24) NOT NULL REFERENCES project_permissions(id),
    PRIMARY KEY (group_id, project_permission_id)
);

CREATE TABLE user_group_memberships (
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    group_id BIGINT NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, group_id)
);

CREATE TABLE global_user_permissions(
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    global_permission_id VARCHAR(24) REFERENCES global_permissions(id),
    PRIMARY KEY (user_id, global_permission_id)
);


-- +goose Down
DROP TABLE global_user_permissions;
DROP TABLE user_group_memberships;
DROP TABLE group_permissions;
DROP TABLE groups;
DROP TABLE project_permissions;
DROP TABLE global_permissions;
