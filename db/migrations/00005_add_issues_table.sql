-- ## Issue
-- ```go
-- id
-- title
-- description // long form description
-- code // serial number per project eg KERR-[[100]]
-- status // todo, progress, done, cancelled. keep it opinionated.
--
-- project_id
-- recruiter_user_id
--
-- priority // an integer. In the UI, will be shown as a name, if a "label" is associated to this specific value.
--          // eg. LOW = 60, NORMAL = 100, HIGH = 150, CRITICAL = 999
--          // with this we can do some crazy calcs
--          // eg. a "heat" statistic which is the average of the priorities
-- ```
-- By the way, in shorthand syntax then use:
-- + for labels 
-- @ for assigning (plus special @NOBODY)
-- ! for predefined priorities, or !<integer> for constant
-- > for dependencies

-- +goose Up
CREATE TYPE status AS ENUM ('todo', 'progress', 'done', 'cancelled');
CREATE TABLE issues (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    code BIGINT NOT NULL,
    description TEXT,

    status STATUS NOT NULL,
    priority INT NOT NULL,

    project_id BIGINT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    recruiter_user_id BIGINT REFERENCES users(id)
);

-- ## IssueAssignees(IssueId, UserId)
CREATE TABLE issue_assignees(
    issue_id BIGINT REFERENCES issues(id) ON DELETE CASCADE,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, issue_id)
);

-- ## IssueLabels(IssueId, LabelId)
CREATE TABLE issue_labels(
    issue_id BIGINT REFERENCES issues(id) ON DELETE CASCADE,
    label_id BIGINT REFERENCES labels(id) ON DELETE CASCADE,
    PRIMARY KEY (issue_id, label_id)
);

-- ## IssueRelationship(FromIssueId, ToIssueId, Kind)
-- Kind = dependency, ...
-- #KERR-12 DEP #KERR-13
-- means 12 depends on 13
--
-- so inbound relationships, check ToIssueId
-- outbounds,                check FromIssueId
CREATE TYPE relationship AS ENUM ('dependency');
CREATE TABLE issue_relationships(
    from_issue_id BIGINT REFERENCES issues(id) ON DELETE CASCADE,
    to_issue_id BIGINT REFERENCES issues(id) ON DELETE CASCADE,
    category RELATIONSHIP NOT NULL,
    PRIMARY KEY (from_issue_id, to_issue_id)
);

-- +goose Down
DROP TABLE issue_relationships;
DROP TYPE relationship;
DROP TABLE issue_assignees;
DROP TABLE issue_labels;
DROP TABLE issues;
DROP TYPE status;
