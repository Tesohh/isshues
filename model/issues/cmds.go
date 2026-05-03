package issues

import (
	"context"
	"errors"
	"slices"

	tea "charm.land/bubbletea/v2"
	"charm.land/log/v2"
	"github.com/Tesohh/isshues/action"
	db "github.com/Tesohh/isshues/db/generated"
	dbmore "github.com/Tesohh/isshues/db/more"
	"github.com/Tesohh/isshues/model"
)

var (
	ErrNotMember = errors.New("cannot load project, as you are not a member")
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
		return model.ErrInternalMsg()
	}
	defer func() { _ = tx.Rollback(ctx) }()
	query := db.New(tx)

	// check permission
	hasPermission, err := query.UserIsMemberOfProject(ctx, db.UserIsMemberOfProjectParams{UserID: m.userId, ProjectID: m.projectId})
	if err != nil {
		log.Error("issues.Model.LoadProjectCmd: error when checking permission", "err", err)
		return model.ErrInternalMsg()
	}

	if !hasPermission {
		return model.ErrMsg{Err: ErrNotMember}
	}

	// get project
	project, err := query.GetProjectById(ctx, m.projectId)
	if err != nil {
		log.Error("issues.Model.LoadProjectCmd: error when querying project", "err", err)
		return model.ErrInternalMsg()
	}

	// load views
	views, err := query.GetAllViewsInProject(ctx, m.projectId)
	if err != nil {
		log.Error("issues.Model.LoadProjectCmd: error when querying views", "err", err)
		return model.ErrInternalMsg()
	}

	return UpdateProjectMsg{Project: project, Views: views}
}

type issueAndRelations struct {
	issue                  db.Issue
	assigneeIDs            []int64
	labelIDs               []int64
	relationshipToIssueIDs []int64
}

type viewData struct {
	issuesMap        map[int64]*issueAndRelations
	issuesOrderedIDs []int64            // maps are by design unordered; which means we have to store the order somewhere..
	users            map[int64]db.User  // list of all assignee users from all issues in this view
	labels           map[int64]db.Label // list of all labels from all issues in this view
	shallowIssues    map[int64]db.Issue // list of all issues with a incoming relationship from all issues in this view
}

type UpdateViewDataMsg struct {
	viewID   int64
	viewData viewData
}

func (m Model) MakeLoadIssuesForSelectedViewCmd() func() tea.Msg {
	id := m.tabs.SelectedID()
	index := slices.IndexFunc(m.views, func(v db.View) bool {
		return v.ID == id
	})
	if index == -1 {
		return nil
	}

	return m.MakeLoadIssuesForViewCmd(m.views[index])
}

// also loads: labels, users (assignees), dependencies.
func (m Model) MakeLoadIssuesForViewCmd(view db.View) func() tea.Msg {
	return func() tea.Msg {
		ctx := context.Background()
		log.Info("loading issues", "userID", m.userId, "view.ID", view.ID, "view.Name", view.Name)

		queryFail := func(err error, name string) model.ErrMsg {
			log.Error("issues.Model.LoadIssueForViewCmd: error when querying", "table", name, "err", err)
			return model.ErrInternalMsg()
		}

		queryStr, binds := dbmore.GenerateViewQuery(view, dbmore.ViewQueryParams{ViewerUserID: m.userId})
		issues, err := action.RunViewQuery(m.app, m.app.DBPool, queryStr, binds)
		if err != nil {
			return queryFail(err, "views")
		}

		issueOrderedIDs := make([]int64, 0, len(issues))
		issueMap := make(map[int64]*issueAndRelations, len(issues))
		for _, issue := range issues {
			issueOrderedIDs = append(issueOrderedIDs, issue.ID)
			issueMap[issue.ID] = &issueAndRelations{issue: issue}
		}

		// load: assignee, label, relationships IDs in bulk (all issues at the same time)
		// put them into the issue map
		// and take their ID for later querying

		issueAssignees, err := m.app.DB.GetIssueAssigneeIDsBulk(ctx, issueOrderedIDs)
		if err != nil {
			return queryFail(err, "issue_assignees")
		}
		issueAssigneeIDs := make([]int64, 0, len(issueAssignees))
		for _, issueAssignee := range issueAssignees {
			target := issueMap[issueAssignee.IssueID]
			target.assigneeIDs = append(target.assigneeIDs, issueAssignee.UserID)
			issueAssigneeIDs = append(issueAssigneeIDs, issueAssignee.UserID)
		}

		issueLabels, err := m.app.DB.GetIssueLabelIDsBulk(ctx, issueOrderedIDs)
		if err != nil {
			return queryFail(err, "issue_labels")
		}
		issueLabelIDs := make([]int64, 0, len(issueLabels))
		for _, issueLabel := range issueLabels {
			target := issueMap[issueLabel.IssueID]
			target.labelIDs = append(target.labelIDs, issueLabel.LabelID)
			issueLabelIDs = append(issueLabelIDs, issueLabel.LabelID)
		}

		issueRelationships, err := m.app.DB.GetIssueRelationshipsBulk(ctx, issueOrderedIDs)
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
		usersList, err := m.app.DB.GetUsersByIDBulk(ctx, issueAssigneeIDs)
		if err != nil {
			return queryFail(err, "users")
		}
		users := make(map[int64]db.User, len(usersList))
		for _, user := range usersList {
			users[user.ID] = user
		}

		labelsList, err := m.app.DB.GetLabelsByIDBulk(ctx, issueLabelIDs)
		if err != nil {
			return queryFail(err, "labels")
		}
		labels := make(map[int64]db.Label, len(labelsList))
		for _, label := range labelsList {
			labels[label.ID] = label
		}

		shallowIssuesList, err := m.app.DB.GetIssuesByIDBulk(ctx, issueRelationshipToIssueIDs)
		if err != nil {
			return queryFail(err, "issues (shallow relationships)")
		}
		shallowIssues := make(map[int64]db.Issue, len(shallowIssuesList))
		for _, shallowIssue := range shallowIssuesList {
			shallowIssues[shallowIssue.ID] = shallowIssue
		}

		return UpdateViewDataMsg{
			viewID: view.ID,
			viewData: viewData{
				issuesMap:        issueMap,
				issuesOrderedIDs: issueOrderedIDs,
				users:            users,
				labels:           labels,
				shallowIssues:    shallowIssues,
			},
		}
	}
}
