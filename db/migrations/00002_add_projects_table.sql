-- +goose Up
-- ## Project(Id, Title, Prefix)
CREATE TABLE projects (
    id BIGSERIAL PRIMARY KEY,
    title text NOT NULL,
    prefix char(4) UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE projects;
