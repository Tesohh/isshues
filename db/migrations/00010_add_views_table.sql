-- +goose Up
CREATE TYPE view_priority_mode AS ENUM ('lt', 'le', 'eq', 'ge', 'gt');
CREATE TYPE view_many_mode AS ENUM ('any', 'all', 'exact');
CREATE TYPE view_sort_by AS ENUM ('code', 'edit_date', 'priority');
CREATE TYPE view_sort_order AS ENUM ('ascending', 'descending');
CREATE TYPE view_style AS ENUM ('panels', 'table');

-- NOTE that users must also have permission to view issues.
CREATE TABLE views (
	id BIGSERIAL PRIMARY KEY,
	project_id BIGINT NOT NULL REFERENCES projects(id),
	name TEXT NOT NULL,

	title TEXT, -- this is given to a LIKE

	statuses status[], -- the "many mode" is `any` in this case

	priority INTEGER,
	priority_mode view_priority_mode NOT NULL DEFAULT 'eq',

	labels_mode view_many_mode NOT NULL DEFAULT 'all',
	
	assignees_mode view_many_mode NOT NULL DEFAULT 'any', -- WARN: maybe 'any' and empty assignees doesn't result in all issues being picked regardless of assignees!
	assignees_include_viewer BOOLEAN NOT NULL DEFAULT false,

	assignee_groups_mode view_many_mode NOT NULL DEFAULT 'any',

	sort_by view_sort_by NOT NULL DEFAULT 'code',
	sort_order view_sort_order NOT NULL DEFAULT 'ascending',

	style view_style NOT NULL DEFAULT 'panels'
);

CREATE TABLE view_labels (
	view_id BIGINT NOT NULL REFERENCES views(id) ON DELETE CASCADE,
	label_id BIGINT NOT NULL REFERENCES labels(id) ON DELETE CASCADE,
	PRIMARY KEY(view_id, label_id)
);

CREATE TABLE view_assignees (
	view_id BIGINT NOT NULL REFERENCES views(id) ON DELETE CASCADE,
	user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	PRIMARY KEY(view_id, user_id)
);

CREATE TABLE view_group_assignees (
	view_id BIGINT NOT NULL REFERENCES views(id) ON DELETE CASCADE,
	group_id BIGINT NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
	PRIMARY KEY(view_id, group_id)
);

-- +goose Down
DROP TABLE view_labels;
DROP TABLE view_assignees;
DROP TABLE view_group_assignees;
DROP TABLE views;
DROP TYPE view_sort_order;
DROP TYPE view_sort_by;
DROP TYPE view_many_mode;
DROP TYPE view_priority_mode;
DROP TYPE view_style;
