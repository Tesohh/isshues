package issuedetail

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	db "github.com/Tesohh/isshues/db/generated"
	"github.com/Tesohh/isshues/ui"
)

func (m Model) MakeHeader() string {
	unreset := func(s string) string {
		return strings.ReplaceAll(s, "\x1b[m", "")
	}

	// unresetMany := func(ss []string) []string {
	// 	for i, s := range ss {
	// 		ss[i] = s[:len(s)-2]
	// 	}
	// 	return ss
	// }

	str := ""
	// HACK reset sequences at the end of the strings cause the background to not be rendered as expected.
	row1 := fmt.Sprintf("%s %s: %s",
		unreset(ui.CompIssueStatusCircle(&m.issue, m.theme)),
		unreset(ui.CompIssueCode(&m.issue, m.theme)),
		unreset(ui.CompIssueTitle(&m.issue, m.theme)),
	)

	str += row1

	// TEMP: handle ok
	recruiter, _ := m.users[m.issue.RecruiterUserID]

	row2 := fmt.Sprintf("by %s", ui.CompIssueRecruiter(&m.issue, m.theme, recruiter))

	if len(m.assigneeIDs) > 0 {
		users := make([]db.User, 0, len(m.assigneeIDs))
		for _, userID := range m.assigneeIDs {
			if user, ok := m.users[userID]; ok {
				users = append(users, user)
			}
		}
		row2 += "-> to " + strings.Join(ui.CompIssueAssignees(&m.issue, m.theme, users, ""), ", ")
	}

	str += "\n" + row2

	row3 := "todo labels"
	if m.issue.Priority != m.app.Viper.GetInt32("priorities.default") {
		row3 += strings.Join(ui.CompIssuePriority(&m.issue, m.theme, m.app.Viper), " ")
	}
	str += "\n" + row3

	bg := lipgloss.NewStyle().
		Background(ui.HLDefs.Get(ui.HLKeySurface, m.theme)).
		Width(m.width)
	return bg.Render(str)
}
