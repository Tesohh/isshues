-- ## Label(Id, Name, Color, ProjectId)
-- In v2 this would correspond to both category roles and tags

-- +goose Up
CREATE TABLE labels (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    symbol CHAR(1), -- single character symbol (from nerd font) shown in some views to save space, eg. +feat -> +󰇈 , +bug -> +
    color CHAR(7), -- color in hex form eg. #AABBCC. If null, will have no background. 
    project_id BIGINT REFERENCES projects(id) ON DELETE CASCADE NOT NULL
);

-- +goose Down
DROP TABLE labels;
