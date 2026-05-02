package config

import db "github.com/Tesohh/isshues/db/generated"

type DefaultView struct {
	Name     string      `mapstructure:"name"`
	Title    string      `mapstructure:"title"`
	Statuses []db.Status `mapstructure:"statuses"`

	Priority     *int                `mapstructure:"priority"`
	PriorityMode db.ViewPriorityMode `mapstructure:"priority_mode"`

	Labels     []string        `mapstructure:"labels"`
	LabelsMode db.ViewManyMode `mapstructure:"labels_mode"`

	Assignees              []string        `mapstructure:"assignees"`
	AssigneesMode          db.ViewManyMode `mapstructure:"assignees_mode"`
	AssigneesIncludeViewer bool            `mapstructure:"assignees_include_viewer"`

	AssigneeGroups     []string        `mapstructure:"assignee_groups"`
	AssigneeGroupsMode db.ViewManyMode `mapstructure:"assignee_groups_mode"`

	SortBy    db.ViewSortBy    `mapstructure:"sort_by"`
	SortOrder db.ViewSortOrder `mapstructure:"sort_order"`

	Style db.ViewStyle `mapstructure:"style"`
}
