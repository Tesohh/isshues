-- +goose Up
CREATE TABLE users_ssh_pks (
    user_id BIGINT NOT NULL REFERENCES users(id),
    ssh_public_key text NOT NULL
);

-- +goose Down
DROP TABLE users_ssh_pks;
