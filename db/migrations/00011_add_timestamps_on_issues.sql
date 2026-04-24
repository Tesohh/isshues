-- +goose Up
BEGIN TRANSACTION;
SET LOCAL lock_timeout = '5s';

ALTER TABLE issues 
ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT now();

-- NOTE: remember to update this!!!
ALTER TABLE issues 
ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT now();

COMMIT;

-- +goose Down
BEGIN TRANSACTION;

ALTER TABLE issues 
    DROP COLUMN IF EXISTS updated_at;

ALTER TABLE issues 
    DROP COLUMN IF EXISTS created_at;

COMMIT TRANSACTION;
