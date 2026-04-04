package action

import (
	"context"

	"github.com/Tesohh/isshues/app"
	db "github.com/Tesohh/isshues/db/generated"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateIssueParams struct {
	Title       string
	Description string
	Priority    int
	RecruiterID int64
	ProjectID   int64

	UserMentionIDs  []int64
	GroupMentionIDs []int64
	DependencyIDs   []int64
	LabelIDs        []int64
}

// 0. gets the amount of issues (for the code)
// 1. inserts the issue
// 2. inserts all assignees
// 3. inserts assignee groups TODO
// 4. inserts label refs
// 5. inserts dependencies
// Assumes the user has the right permissions, and all input is valid
func CreateIssue(app *app.App, params CreateIssueParams) (db.Issue, error) {
	// TODO: use a db transaction

	ctx := context.Background()

	// 0. gets the amount of issues (for the code)
	count, err := app.DB.GetIssuesCountInProject(ctx, params.ProjectID)
	if err != nil {
		return db.Issue{}, err
	}

	// 1. inserts the issue
	issue, err := app.DB.InsertIssue(ctx, db.InsertIssueParams{
		Title:           params.Title,
		Code:            count + 1,
		Description:     pgtype.Text{String: params.Description, Valid: true},
		Status:          db.StatusTodo,
		Priority:        int32(params.Priority),
		ProjectID:       params.ProjectID,
		RecruiterUserID: params.RecruiterID,
	})
	if err != nil {
		return issue, err
	}

	// 2. inserts all assignees
	assigneeParams := []db.BulkInsertIssueAssigneesParams{}
	for _, id := range params.UserMentionIDs {
		assigneeParams = append(assigneeParams, db.BulkInsertIssueAssigneesParams{IssueID: issue.ID, UserID: id})
	}

	_, err = app.DB.BulkInsertIssueAssignees(ctx, assigneeParams)
	if err != nil {
		return issue, err
	}

	// 3. inserts assignee groups TODO

	// 4. inserts label refs
	labelParams := []db.BulkInsertIssueLabelsParams{}
	for _, id := range params.LabelIDs {
		labelParams = append(labelParams, db.BulkInsertIssueLabelsParams{IssueID: issue.ID, LabelID: id})
	}

	_, err = app.DB.BulkInsertIssueLabels(ctx, labelParams)
	if err != nil {
		return issue, err
	}

	// 5. inserts dependencies
	relationshipParams := make([]db.BulkInsertIssueRelationshipsParams, len(params.DependencyIDs))
	for _, id := range params.DependencyIDs {
		relationshipParams = append(relationshipParams, db.BulkInsertIssueRelationshipsParams{FromIssueID: issue.ID, ToIssueID: id, Category: db.RelationshipDependency})
	}

	_, err = app.DB.BulkInsertIssueLabels(ctx, labelParams)
	if err != nil {
		return issue, err
	}

	return issue, nil
}
