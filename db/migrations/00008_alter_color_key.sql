-- +goose Up
ALTER TABLE groups ALTER COLUMN color TYPE VARCHAR(16);
ALTER TABLE groups RENAME COLUMN color TO color_key; 

ALTER TABLE labels ALTER COLUMN color TYPE VARCHAR(16);
ALTER TABLE labels RENAME COLUMN color TO color_key; 

-- +goose Down
ALTER TABLE groups RENAME COLUMN color_key TO color;
ALTER TABLE groups ALTER COLUMN color TYPE CHAR(7);
ALTER TABLE labels RENAME COLUMN color_key TO color;
ALTER TABLE labels ALTER COLUMN color TYPE CHAR(7);
