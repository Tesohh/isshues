package issues

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/list"
	"github.com/Tesohh/isshues/app"
	db "github.com/Tesohh/isshues/db/generated"
	tint "github.com/lrstanley/bubbletint/v2"
)

// item provides the issue's prerendered title and description.
// delegates should not rerender title and description
type item struct {
	plainTitle  string
	title, desc string
}

func itemFromIssue(app *app.App, theme *tint.Tint, i *issueAndRelations, viewData *viewData) item {
	users := make([]db.User, 0, len(i.assigneeIDs))
	for _, id := range i.assigneeIDs {
		if user, ok := viewData.users[id]; ok {
			users = append(users, user)
		}
	}

	labels := make([]db.Label, 0, len(i.labelIDs))
	for _, id := range i.labelIDs {
		if label, ok := viewData.labels[id]; ok {
			labels = append(labels, label)
		}
	}

	shallowIssues := make([]db.Issue, 0, len(i.relationshipToIssueIDs))
	for _, id := range i.relationshipToIssueIDs {
		if shallowIssue, ok := viewData.shallowIssues[id]; ok {
			shallowIssues = append(shallowIssues, shallowIssue)
		}
	}

	title := fmt.Sprintf("%s %s: %s %s",
		ComponentStatusCircle(&i.issue, theme),
		ComponentCode(&i.issue, theme),
		ComponentTitle(&i.issue, theme),
		ComponentDescription(&i.issue, theme),
	)

	bottomStrs := []string{}
	bottomStrs = append(bottomStrs, ComponentPriority(&i.issue, theme, app.Viper)...)
	bottomStrs = append(bottomStrs, ComponentLabels(&i.issue, theme, labels)...)
	bottomStrs = append(bottomStrs, ComponentDependencies(&i.issue, theme, shallowIssues)...)
	bottomStrs = append(bottomStrs, ComponentAssignees(&i.issue, theme, users, "")...) // TODO: session.User()
	// TODO: handle groups

	desc := strings.Join(bottomStrs, " ")

	return item{
		plainTitle: fmt.Sprintf("#%d: %s", i.issue.Code, i.issue.Title),
		title:      title,
		desc:       desc,
	}
}

func itemsFromViewData(app *app.App, theme *tint.Tint, viewData *viewData) []list.Item {
	items := []list.Item{}
	for _, id := range viewData.issuesOrderedIDs {
		issueAndRelations := viewData.issuesMap[id]
		items = append(items, itemFromIssue(app, theme, issueAndRelations, viewData))
	}
	return items
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.plainTitle }
