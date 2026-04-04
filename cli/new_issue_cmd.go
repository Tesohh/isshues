package cli

import (
	"context"
	"errors"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/log/v2"
	"github.com/Tesohh/isshues/action"
	"github.com/Tesohh/isshues/app"
	db "github.com/Tesohh/isshues/db/generated"
	"github.com/Tesohh/isshues/shorthand"
	"github.com/charmbracelet/ssh"
	"github.com/jackc/pgx/v5"
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
				Title:           product.Text,
				Description:     "",
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
			// broadcast a RefreshIssues{this project} to everyone

			return nil
		},
	}

	return newCmd
}
