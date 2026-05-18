package ui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/Tesohh/isshues/config"
	db "github.com/Tesohh/isshues/db/generated"
	tint "github.com/lrstanley/bubbletint/v2"
	"github.com/spf13/viper"
)

func CompIssueStatusCircle(issue *db.Issue, theme *tint.Tint) string {
	color := HLDefs.Get(HLKey("status-"+issue.Status), theme)
	return lipgloss.NewStyle().Foreground(color).Render("◉")
}

func CompIssueCode(issue *db.Issue, theme *tint.Tint) string {
	mutedStyle := lipgloss.NewStyle().Foreground(HLDefs.Get(HLKeyMuted, theme))

	return mutedStyle.Render(fmt.Sprintf("#%d", issue.Code))
}

func CompIssuePrefixAndCode(issue *db.Issue, project *db.Project, theme *tint.Tint) string {
	mutedStyle := lipgloss.NewStyle().Foreground(HLDefs.Get(HLKeyMuted, theme))

	return mutedStyle.Render(fmt.Sprintf("#%s-%d", project.Prefix, issue.Code))
}

func CompIssueTitle(issue *db.Issue, theme *tint.Tint) string {
	return lipgloss.NewStyle().Foreground(HLDefs.Get(HLKeyText, theme)).Render(issue.Title)
}

func CompIssueDescription(issue *db.Issue, theme *tint.Tint) string {
	mutedStyle := lipgloss.NewStyle().Foreground(HLDefs.Get(HLKeyMuted, theme))

	if issue.Description.Valid {
		return mutedStyle.Render("[...]")
	}

	return ""
}

func CompIssuePriority(issue *db.Issue, theme *tint.Tint, viper *viper.Viper) []string {
	var priorities config.Priorities
	_ = viper.UnmarshalKey("priorities", &priorities) // WARN ignored error, might be a problem?
	closestPriority, closestPriorityK := priorities.FindClosest(int(issue.Priority))

	var (
		priorityStyle = lipgloss.NewStyle().Foreground(KeyToColor(theme, closestPriority.ColorKey))
	)

	if closestPriorityK != "default" {
		text := ""
		if closestPriority.Value == int(issue.Priority) {
			text = closestPriorityK
		} else {
			text = fmt.Sprint(issue.Priority)
		}

		return []string{priorityStyle.Render("!" + text)}
	}
	return []string{}
}

func CompIssueLabels(_ *db.Issue, theme *tint.Tint, labels []db.Label) []string {
	mutedColor := HLDefs.Get(HLKeyMuted, theme)
	strs := []string{}
	for _, label := range labels {
		style := lipgloss.NewStyle().Foreground(NullableKeyToColor(theme, mutedColor, label.ColorKey))

		if label.Symbol.Valid {
			strs = append(strs, style.Render("+"+label.Symbol.String+" "))
		} else {
			strs = append(strs, style.Render("+"+label.Name))
		}
	}
	return strs
}

func CompIssueDependencies(_ *db.Issue, theme *tint.Tint, dependencies []db.Issue) []string {
	mutedStyle := lipgloss.NewStyle().Foreground(HLDefs.Get(HLKeyMuted, theme))

	strs := make([]string, 0, len(dependencies))
	for _, dep := range dependencies {
		strs = append(strs, mutedStyle.Render(fmt.Sprintf(">%d", dep.Code)))
	}

	return strs
}

func CompIssueAssignees(_ *db.Issue, theme *tint.Tint, assignees []db.User, thisUsername string) []string {
	mutedStyle := lipgloss.NewStyle().Foreground(HLDefs.Get(HLKeyMuted, theme))
	emphStyle := lipgloss.NewStyle().Foreground(HLDefs.Get(HLKeyEmphasis, theme))

	strs := []string{}
	for _, assignee := range assignees {
		var style *lipgloss.Style

		if assignee.Username == strings.ToLower(thisUsername) {
			style = &emphStyle
		} else {
			style = &mutedStyle
		}

		if assignee.Shortname.Valid {
			strs = append(strs, style.Render(fmt.Sprintf("@%s", assignee.Shortname.String)))
		} else {
			strs = append(strs, style.Render(fmt.Sprintf("@%s", assignee.Username)))
		}
	}

	return strs
}

func CompIssueRecruiter(_ *db.Issue, theme *tint.Tint, user db.User) string {
	return lipgloss.NewStyle().Foreground(HLDefs.Get(HLKeyText, theme)).Render("@" + user.Username)
}
