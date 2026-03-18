-- +goose Up
CREATE TABLE user_ssh_keys (
    user_id BIGINT NOT NULL REFERENCES users(id),
    ssh_public_key text NOT NULL
);

-- +goose Down
DROP TABLE user_ssh_keys;
