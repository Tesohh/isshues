package issues

import (
	"context"
	"errors"
	"maps"
	"slices"

	tea "charm.land/bubbletea/v2"
	"charm.land/log/v2"
	"github.com/Tesohh/isshues/action"
	db_complex "github.com/Tesohh/isshues/db/complex"
	db "github.com/Tesohh/isshues/db/generated"
	"github.com/Tesohh/isshues/model"
)

var (
	NotMemberErr = errors.New("cannot load project, as you are not a member")
)

type UpdateProjectMsg struct {
	Project db.Project
	Views   []db.View
}

// also loads views, as you never want to change the project and not the views.
func (m Model) LoadProjectCmd() tea.Msg {
	ctx := context.Background()
	// start transaction
	tx, err := m.app.DBPool.Begin(ctx)
	if err != nil {
		log.Error("issues.Model.LoadProjectCmd: cannot start transaction", "err", err)
		return model.InternalErrMsg()
	}
	defer tx.Rollback(ctx)
	query := db.New(tx)

	// check permission
	hasPermission, err := query.UserIsMemberOfProject(ctx, db.UserIsMemberOfProjectParams{UserID: m.userId, ProjectID: m.projectId})
	if err != nil {
		log.Error("issues.Model.LoadProjectCmd: error when checking permission", "err", err)
		return model.InternalErrMsg()
	}

	if !hasPermission {
		return model.ErrMsg{Err: NotMemberErr}
	}

	// get project
	project, err := query.GetProjectById(ctx, m.projectId)
	if err != nil {
		log.Error("issues.Model.LoadProjectCmd: error when querying project", "err", err)
		return model.InternalErrMsg()
	}

	// load views
	views, err := query.GetAllViewsInProject(ctx, m.projectId)
	if err != nil {
		log.Error("issues.Model.LoadProjectCmd: error when querying views", "err", err)
		return model.InternalErrMsg()
	}

	return UpdateProjectMsg{Project: project, Views: views}
}

type issueAndRelations struct {
	issue                  db.Issue
	assigneeIDs            []int64
	labelIDs               []int64
	relationshipToIssueIDs []int64
}

type UpdateViewDataMsg struct {
	viewID   int64
	viewData viewData
}

// also loads: labels, users (assignees), dependencies.
func (m Model) MakeLoadIssuesForViewCmd(view db.View) func() tea.Msg {
	return func() tea.Msg {
		ctx := context.Background()

		queryFail := func(err error, name string) model.ErrMsg {
			log.Error("issues.Model.LoadIssueForViewCmd: error when querying", "table", name, "err", err)
			return model.InternalErrMsg()
		}

		// tx, err := m.app.DBPool.Begin(ctx)
		// if err != nil {
		// 	log.Error("issues.Model.LoadProjectCmd: cannot start transaction", "err", err)
		// 	return model.InternalErrMsg()
		// }
		// defer tx.Rollback(ctx)

		queryStr, binds := db_complex.GenerateViewQuery(view, db_complex.ViewQueryParams{ViewerUserID: m.userId})
		issues, err := action.RunViewQuery(m.app, m.app.DBPool, queryStr, binds)
		if err != nil {
			return queryFail(err, "views")
		}

		issueMap := make(map[int64]*issueAndRelations, len(issues))
		for _, issue := range issues {
			issueMap[issue.ID] = &issueAndRelations{issue: issue}
		}
		issueIDs := slices.Collect(maps.Keys(issueMap))

		// load: assignee, label, relationships IDs in bulk (all issues at the same time)
		// put them into the issue map
		// and take their ID for later querying

		issueAssignees, err := m.app.DB.GetIssueAssigneeIDsBulk(ctx, issueIDs)
		if err != nil {
			return queryFail(err, "issue_assignees")
		}
		issueAssigneeIDs := make([]int64, 0, len(issueAssignees))
		for _, issueAssignee := range issueAssignees {
			target := issueMap[issueAssignee.IssueID]
			target.assigneeIDs = append(target.assigneeIDs, issueAssignee.UserID)
			issueAssigneeIDs = append(issueAssigneeIDs, issueAssignee.UserID)
		}

		issueLabels, err := m.app.DB.GetIssueLabelIDsBulk(ctx, issueIDs)
		if err != nil {
			return queryFail(err, "issue_labels")
		}
		issueLabelIDs := make([]int64, 0, len(issueLabels))
		for _, issueLabel := range issueLabels {
			target := issueMap[issueLabel.IssueID]
			target.labelIDs = append(target.labelIDs, issueLabel.LabelID)
			issueLabelIDs = append(issueLabelIDs, issueLabel.LabelID)
		}

		issueRelationships, err := m.app.DB.GetIssueRelationshipsBulk(ctx, issueIDs)
		if err != nil {
			return queryFail(err, "issue_relationships")
		}
		issueRelationshipToIssueIDs := make([]int64, 0, len(issueRelationships))
		for _, issueRelationship := range issueRelationships {
			target := issueMap[issueRelationship.FromIssueID]
			target.relationshipToIssueIDs = append(target.relationshipToIssueIDs, issueRelationship.ToIssueID)
			issueRelationshipToIssueIDs = append(issueRelationshipToIssueIDs, issueRelationship.ToIssueID)
		}

		// load: assignee users, labels, relationship issues (shallowly)
		users, err := m.app.DB.GetUsersByIDBulk(ctx, issueAssigneeIDs)
		if err != nil {
			return queryFail(err, "users")
		}
		labels, err := m.app.DB.GetLabelsByIDBulk(ctx, issueLabelIDs)
		if err != nil {
			return queryFail(err, "labels")
		}
		shallowIssues, err := m.app.DB.GetIssuesByIDBulk(ctx, issueRelationshipToIssueIDs)
		if err != nil {
			return queryFail(err, "issues (shallow relationships)")
		}

		return UpdateViewDataMsg{
			viewID: view.ID,
			viewData: viewData{
				issuesMap:     issueMap,
				users:         users,
				labels:        labels,
				shallowIssues: shallowIssues,
			},
		}
	}
}
