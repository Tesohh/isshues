-- +goose Up
CREATE TABLE user_ssh_keys (
    user_id BIGINT NOT NULL REFERENCES users(id),
    fingerprint CHAR(44) NOT NULL
);

-- +goose Down
DROP TABLE user_ssh_keys;
