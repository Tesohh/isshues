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
	db "github.com/Tesohh/isshues/db/generated"
	"github.com/Tesohh/isshues/model/issues"
	"github.com/Tesohh/isshues/shorthand"
	"github.com/charmbracelet/ssh"
	"github.com/jackc/pgx/v5"
	tint "github.com/lrstanley/bubbletint/v2"
	"github.com/spf13/cobra"
)

var (
	ErrProjectNotFound  = errors.New("project does not exist")
	ErrPermissionDenied = errors.New("you are missing a permission")
)

func newIssueCmd(session ssh.Session, app *app.App, _ **tea.Program) *cobra.Command {
	newCmd := &cobra.Command{
		Use:  "new [project prefix] [shorthand...]",
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			userId, ok := app.SessionIdToUserIds[session.Context().SessionID()]
			if !ok {
				return errors.New("your userid was not found in the session map. might be an auth issue")
			}

			tx, err := app.DBPool.Begin(ctx)
			if err != nil {
				log.Error("new: cannot start transaction", "prefix", args[0], "err", err)
				return ErrInternal
			}
			defer func() { _ = tx.Rollback(ctx) }()

			query := db.New(tx)

			// get the project id
			project, err := query.GetProjectByPrefix(ctx, strings.ToUpper(args[0]))
			if err == pgx.ErrNoRows {
				return ErrProjectNotFound
			} else if err != nil {
				log.Error("new: error while fetching project from prefix", "prefix", args[0], "err", err)
				return ErrInternal
			}

			// check permissions
			hasPermission, err := query.UserHasProjectPermission(ctx, db.UserHasProjectPermissionParams{
				UserID:              userId,
				ProjectPermissionID: "write-issues",
				ProjectID:           project.ID,
			})
			if err != nil && err != pgx.ErrNoRows {
				log.Error("new: error while checking permissions", "err", err)
				return ErrInternal
			}

			if !hasPermission {
				return fmt.Errorf("%w: write-issues", ErrPermissionDenied)
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
				return ErrInternal
			}

			// add the pending labels and collect all ids
			newLabels, err := action.BulkInsertLabels(app, project.ID, product.PendingLabels)
			if err != nil {
				log.Error("new: error while inserting new labels", "err", err)
				return ErrInternal
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
			issue, err := action.CreateIssue(app, query, params)
			if err != nil {
				log.Error("new: error while creating issue", "err", err)
				return ErrInternal
			}

			// load user theme
			settings, err := app.DB.GetUserSettings(ctx, userId)
			if err != nil {
				log.Error("settings query error", "err", err, "userId", userId)
				return ErrInternal
			}

			theme, ok := tint.GetTint(settings.Theme)
			if !ok {
				return fmt.Errorf("%w: %s. go to https://lrstanley.github.io/bubbletint/ to find a list of supported themes", ErrThemeNotFound, settings.Theme)
			}

			// show feedback and warnings

			title := fmt.Sprintf("%s %s: %s %s",
				issues.ComponentStatusCircle(&issue, theme),
				issues.ComponentPrefixAndCode(&issue, &project, theme),
				issues.ComponentTitle(&issue, theme),
				issues.ComponentDescription(&issue, theme),
			)

			bottomStrs := []string{}
			bottomStrs = append(bottomStrs, issues.ComponentPriority(&issue, theme, app.Viper)...)
			bottomStrs = append(bottomStrs, issues.ComponentLabels(&issue, theme, allLabels)...)
			bottomStrs = append(bottomStrs, issues.ComponentDependencies(&issue, theme, product.Dependencies)...)
			bottomStrs = append(bottomStrs, issues.ComponentAssignees(&issue, theme, product.UserMentions, session.User())...)
			// TODO: handle groups

			bottom := strings.Join(bottomStrs, " ")

			warningBgStyle := lipgloss.NewStyle().Background(theme.Yellow).Foreground(theme.Bg)
			warningFgStyle := lipgloss.NewStyle().Foreground(theme.Yellow)
			warningStrs := []string{}
			for _, warning := range product.Warnings {
				plate := warningBgStyle.Render(" WARNING ")
				text := warningFgStyle.Render(warning.Error())
				warningStrs = append(warningStrs, fmt.Sprintf("%s %s", plate, text))
			}

			warnings := strings.Join(warningStrs, "\n")
			if len(warningStrs) > 0 {
				warnings += "\n\n"
			}

			feedback := fmt.Sprintf("Created new issue successfully\n\n%s\n%s\n\n%s", title, bottom, warnings)

			wish.Print(session, feedback)

			err = tx.Commit(ctx)
			if err != nil {
				log.Errorf("new: %s", err.Error())
				return ErrInternal
			}

			// broadcast a RefreshIssues{this project} to everyone

			return nil
		},
	}

	return newCmd
}
