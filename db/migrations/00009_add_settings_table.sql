-- +goose Up
CREATE TABLE user_settings (
    user_id BIGSERIAL PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    theme text NOT NULL
);

-- +goose Down
DROP TABLE user_settings;
