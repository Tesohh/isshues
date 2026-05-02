package issues

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/Tesohh/isshues/common"
	"github.com/Tesohh/isshues/config"
	db "github.com/Tesohh/isshues/db/generated"
	tint "github.com/lrstanley/bubbletint/v2"
	"github.com/spf13/viper"
)

const ComponentMutedDarkenFactor = 0.4
const ComponentEmphDarkenFactor = 0.2

func ComponentStatusCircle(issue *db.Issue, theme *tint.Tint) string {
	var color *tint.Color

	switch issue.Status {
	case db.StatusTodo:
		color = theme.Green
	case db.StatusProgress:
		color = theme.Blue
	case db.StatusDone:
		color = theme.Purple
	case db.StatusCancelled:
		color = theme.Red
	default:
		color = theme.Fg
	}

	return lipgloss.NewStyle().Foreground(color).Render("◉")
}

func ComponentCode(issue *db.Issue, theme *tint.Tint) string {
	mutedColor := lipgloss.Darken(theme.Fg, ComponentMutedDarkenFactor)
	mutedStyle := lipgloss.NewStyle().Foreground(mutedColor)

	return mutedStyle.Render(fmt.Sprintf("#%d", issue.Code))
}

func ComponentPrefixAndCode(issue *db.Issue, project *db.Project, theme *tint.Tint) string {
	mutedColor := lipgloss.Darken(theme.Fg, ComponentMutedDarkenFactor)
	mutedStyle := lipgloss.NewStyle().Foreground(mutedColor)

	return mutedStyle.Render(fmt.Sprintf("#%s-%d", project.Prefix, issue.Code))
}

func ComponentTitle(issue *db.Issue, theme *tint.Tint) string {
	return lipgloss.NewStyle().Foreground(theme.Fg).Render(issue.Title)
}

func ComponentPriority(issue *db.Issue, theme *tint.Tint, viper *viper.Viper) []string {
	var priorities config.Priorities
	viper.UnmarshalKey("priorities", &priorities)
	closestPriority, closestPriorityK := priorities.FindClosest(int(issue.Priority))

	var (
		priorityStyle = lipgloss.NewStyle().Foreground(common.KeyToColor(theme, closestPriority.ColorKey))
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

func ComponentLabels(_ *db.Issue, theme *tint.Tint, labels []db.Label) []string {
	mutedColor := lipgloss.Darken(theme.Fg, ComponentMutedDarkenFactor)
	strs := []string{}
	for _, label := range labels {
		style := lipgloss.NewStyle().Foreground(common.NullableKeyToColor(theme, mutedColor, label.ColorKey))

		if label.Symbol.Valid {
			strs = append(strs, style.Render("+"+label.Symbol.String+" "))
		} else {
			strs = append(strs, style.Render("+"+label.Name))
		}
	}
	return strs
}

func ComponentDependencies(_ *db.Issue, theme *tint.Tint, dependencies []db.Issue) []string {
	mutedColor := lipgloss.Darken(theme.Fg, ComponentMutedDarkenFactor)
	mutedStyle := lipgloss.NewStyle().Foreground(mutedColor)

	strs := make([]string, 0, len(dependencies))
	for _, dep := range dependencies {
		strs = append(strs, mutedStyle.Render(fmt.Sprintf(">%d", dep.Code)))
	}

	return strs
}

func ComponentAssignees(_ *db.Issue, theme *tint.Tint, assignees []db.User, thisUsername string) []string {
	mutedColor := lipgloss.Darken(theme.Fg, ComponentMutedDarkenFactor)
	mutedStyle := lipgloss.NewStyle().Foreground(mutedColor)
	emphStyle := lipgloss.NewStyle().Foreground(lipgloss.Darken(theme.Purple, ComponentEmphDarkenFactor))

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
