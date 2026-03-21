-- +goose Up
CREATE TABLE users(
    id BIGSERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    shortname VARCHAR(4) UNIQUE -- a short alias, eg. @lallos -> @ls, @tesohh -> @ts, @mat1232123 -> @mt
    -- is_admin BOOL NOT NULL
);

-- +goose Down
DROP TABLE users;
