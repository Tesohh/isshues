package cli

import (
	"context"
	"errors"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/log/v2"
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

			// prepare the issue

			// add issue to the db
			// TODO: create pending labels
			// show feedback
			// broadcast a RefreshIssues{this project} to everyone

			return nil
		},
	}

	return newCmd
}
