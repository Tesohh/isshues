package cli

import (
	"context"
	"errors"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"charm.land/log/v2"
	"charm.land/wish/v2"
	"github.com/Tesohh/isshues/action"
	"github.com/Tesohh/isshues/app"
	"github.com/Tesohh/isshues/common"
	"github.com/Tesohh/isshues/config"
	db "github.com/Tesohh/isshues/db/generated"
	"github.com/Tesohh/isshues/shorthand"
	"github.com/charmbracelet/ssh"
	"github.com/jackc/pgx/v5"
	tint "github.com/lrstanley/bubbletint/v2"
	"github.com/spf13/cobra"
)

var (
	ProjectNotFoundErr  = errors.New("project does not exist")
	PermissionDeniedErr = errors.New("you are missing a permission")
)

func newIssueCmd(session ssh.Session, app *app.App, _ **tea.Program) *cobra.Command {
	newCmd := &cobra.Command{
		Use:  "new [project prefix] [shorthand...]",
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			userId, ok := app.SessionIdToUserIds[session.Context().SessionID()]
			if !ok {
				return errors.New("your userid was not found in the session map. might be an auth issue.")
			}

			// get the project id
			project, err := app.DB.GetProjectByPrefix(ctx, strings.ToUpper(args[0]))
			if err == pgx.ErrNoRows {
				return ProjectNotFoundErr
			} else if err != nil {
				log.Error("new: error while fetching project from prefix", "prefix", args[0], "err", err)
				return InternalErr
			}

			// check permissions
			hasPermission, err := app.DB.UserHasProjectPermission(ctx, db.UserHasProjectPermissionParams{
				UserID:              userId,
				ProjectPermissionID: "write-issues",
				ProjectID:           project.ID,
			})
			if err != nil && err != pgx.ErrNoRows {
				log.Error("new: error while checking permissions", "err", err)
				return InternalErr
			}

			if !hasPermission {
				return fmt.Errorf("%w: write-issues", PermissionDeniedErr)
			}

			if len(args) == 1 {
				// TODO: launch form
				return nil
			}

			// if args > 1:
			// merge text
			msg := strings.Join(args[1:], " ")

			// parse the shorthand
			captures := shorthand.Parse(msg)
			product, err := shorthand.Process(captures, app, project.ID, userId)
			if err != nil {
				log.Error("new: error while processing issue", "err", err, "captures", captures)
				return InternalErr
			}

			// add the pending labels and collect all ids
			newLabels, err := action.BulkInsertLabels(app, project.ID, product.PendingLabels)
			if err != nil {
				log.Error("new: error while inserting new labels", "err", err)
				return InternalErr
			}

			// collect ids (if only we had .iter().map()...)
			allLabels := append(newLabels, product.Labels...)
			labelIds := make([]int64, 0, len(allLabels))
			for _, label := range allLabels {
				labelIds = append(labelIds, label.ID)
			}

			userMentionIDs := make([]int64, 0, len(product.UserMentions))
			for _, mention := range product.UserMentions {
				userMentionIDs = append(userMentionIDs, mention.ID)
			}

			groupMentionIDs := make([]int64, 0, len(product.GroupMentions))
			for _, mention := range product.GroupMentions {
				groupMentionIDs = append(groupMentionIDs, mention.ID)
			}

			dependencyIDs := make([]int64, 0, len(product.Dependencies))
			for _, dependency := range product.Dependencies {
				dependencyIDs = append(dependencyIDs, dependency.ID)
			}

			// prepare the issue
			params := action.CreateIssueParams{
				Title:           product.Title,
				Description:     product.Description,
				Priority:        product.Priority,
				RecruiterID:     userId,
				ProjectID:       project.ID,
				UserMentionIDs:  userMentionIDs,
				GroupMentionIDs: groupMentionIDs,
				DependencyIDs:   dependencyIDs,
				LabelIDs:        labelIds,
			}

			log.Info("adding issue with params", "params", params)

			// add issue to the db
			issue, err := action.CreateIssue(app, params)
			_ = issue
			if err != nil {
				log.Error("new: error while creating issue", "err", err)
				return InternalErr
			}

			// show feedback and warnings
			// TODO: put this in a reusable function
			// TODO: use requested theme

			tint := tint.TintRosePine

			var priorities config.Priorities
			app.Viper.UnmarshalKey("priorities", &priorities)
			closestPriority, closestPriorityK := priorities.FindClosest(int(issue.Priority))

			var (
				mutedColor     = lipgloss.Darken(tint.Fg, 0.4)
				mutedStyle     = lipgloss.NewStyle().Foreground(mutedColor)
				emphStyle      = lipgloss.NewStyle().Foreground(lipgloss.Darken(tint.Purple, 0.2))
				priorityStyle  = lipgloss.NewStyle().Foreground(common.KeyToColor(tint, closestPriority.ColorKey))
				warningBgStyle = lipgloss.NewStyle().Background(tint.Yellow).Foreground(tint.Bg)
				warningFgStyle = lipgloss.NewStyle().Foreground(tint.Yellow)
			)

			circleStr := common.MakeStatusCircle(tint, issue.Status)
			serialStr := mutedStyle.Render(fmt.Sprintf("#%s-%d:", project.Prefix, issue.Code))
			title := fmt.Sprintf("%s %s %s", circleStr, serialStr, issue.Title)

			if issue.Description.Valid {
				title += " " + mutedStyle.Render("[...]")
			}

			bottomStrs := []string{}

			if closestPriorityK != "default" {
				text := ""
				if closestPriority.Value == int(issue.Priority) {
					text = closestPriorityK
				} else {
					text = fmt.Sprint(issue.Priority)
				}

				bottomStrs = append(bottomStrs, priorityStyle.Render("!"+text))
			}

			for _, label := range allLabels {
				style := lipgloss.NewStyle().Foreground(common.NullableKeyToColor(tint, mutedColor, label.ColorKey))

				if label.Symbol.Valid {
					bottomStrs = append(bottomStrs, style.Render("+"+label.Symbol.String+" "))
				} else {
					bottomStrs = append(bottomStrs, style.Render("+"+label.Name))
				}
			}

			for _, dep := range product.Dependencies {
				bottomStrs = append(bottomStrs, mutedStyle.Render(fmt.Sprintf(">%d", dep.Code)))
			}

			// TODO: handle groups
			for _, assignee := range product.UserMentions {
				var style *lipgloss.Style

				if assignee.Username == strings.ToLower(session.User()) {
					style = &emphStyle
				} else {
					style = &mutedStyle
				}

				if assignee.Shortname.Valid {
					bottomStrs = append(bottomStrs, style.Render(fmt.Sprintf("@%s", assignee.Shortname.String)))
				} else {
					bottomStrs = append(bottomStrs, style.Render(fmt.Sprintf("@%s", assignee.Username)))
				}
			}

			bottom := strings.Join(bottomStrs, " ")

			// TODO show warnings
			warningStrs := []string{}
			for _, warning := range product.Warnings {
				plate := warningBgStyle.Render(" WARNING ")
				text := warningFgStyle.Render(warning.Error())
				warningStrs = append(warningStrs, fmt.Sprintf("%s: %s", plate, text))
			}

			warnings := strings.Join(warningStrs, "\n")
			if len(warningStrs) > 0 {
				warnings += "\n\n"
			}

			feedback := fmt.Sprintf("Created new issue successfully\n\n%s\n%s\n\n%s", title, bottom, warnings)

			wish.Print(session, feedback)

			// broadcast a RefreshIssues{this project} to everyone

			return nil
		},
	}

	return newCmd
}
