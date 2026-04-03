package action

import (
	"github.com/Tesohh/isshues/app"
	db "github.com/Tesohh/isshues/db/generated"
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

// 1. inserts the issue
// 2. inserts all assignees
// 3. inserts assignee groups TODO
// 4. inserts label refs
// 5. inserts dependencies
// Assumes the user has the right permissions, and all input is valid
func CreateIssue(app *app.App, params CreateIssueParams) (*db.Issue, error)
