package issuedetail

import (
	db "github.com/Tesohh/isshues/db/generated"
	tint "github.com/lrstanley/bubbletint/v2"
)

type Model struct {
	issue db.Issue

	assigneeIDs            []int64 // all assignees referenced by issue
	labelIDs               []int64 // all labels referenced by issue
	relationshipToIssueIDs []int64 // all issues with a incoming relationship from issue

	users         map[int64]db.User  // all users existing in the view. naturally, it must include at least all ids in assigneeIDs
	labels        map[int64]db.Label // all labels existing in the view. naturally, it must include at least all ids in labelIDs
	shallowIssues map[int64]db.Issue // all shallowIssues existing in the view. naturally, it must include at least all ids in relationshipToIssueIDs

	width, height int

	theme *tint.Tint
}

func New() Model {
	return Model{
		issue:                  db.Issue{},
		assigneeIDs:            []int64{},
		labelIDs:               []int64{},
		relationshipToIssueIDs: []int64{},
		users:                  map[int64]db.User{},
		labels:                 map[int64]db.Label{},
		shallowIssues:          map[int64]db.Issue{},
		width:                  0,
		height:                 0,
		theme:                  &tint.Tint{},
	}
}

func (m Model) SetTheme(theme *tint.Tint) Model {
	m.theme = theme
	return m
}

func (m Model) SetSize(width, height int) Model {
	m.height = height
	m.width = width
	return m
}

func (m Model) SetIssueData(issue db.Issue, assigneeIDs, labelIDs, relationshipToIssueIDs []int64) Model {
	m.issue = issue
	m.assigneeIDs = assigneeIDs
	m.labelIDs = labelIDs
	m.relationshipToIssueIDs = relationshipToIssueIDs
	return m
}

func (m Model) SetViewData(users map[int64]db.User, labels map[int64]db.Label, shallowIssues map[int64]db.Issue) Model {
	m.users = users
	m.labels = labels
	m.shallowIssues = shallowIssues
	return m
}
