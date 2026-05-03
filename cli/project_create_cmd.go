package cli

import (
	"context"
	"errors"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/log/v2"
	"github.com/Tesohh/isshues/action"
	"github.com/Tesohh/isshues/app"
	db "github.com/Tesohh/isshues/db/generated"
	"github.com/Tesohh/isshues/model/projects"
	"github.com/charmbracelet/ssh"
	"github.com/spf13/cobra"
)

func projectCreateCmd(session ssh.Session, app *app.App, _ **tea.Program) *cobra.Command {
	createCmd := &cobra.Command{
		Use:   "create [prefix] [title]",
		Args:  cobra.MinimumNArgs(2),
		Short: "Creates a new project (requires the `create-projects` permission)",
		Long:  "Creates a new project with the given prefix, and title;\nthen creates default groups and adds you to the add_creator groups specified in the server config\n(requires the `create-projects` permission)",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			tx, err := app.DBPool.Begin(ctx)
			if err != nil {
				log.Error("project create: cannot start transaction", "prefix", args[0], "err", err)
				return ErrInternal
			}
			defer func() { _ = tx.Rollback(ctx) }()

			query := db.New(tx)

			userId, ok := app.SessionIdToUserIds[session.Context().SessionID()]
			if !ok {
				return errors.New("your userid was not found in the session map. might be an auth issue")
			}

			authorized, err := query.UserHasGlobalPermission(ctx, db.UserHasGlobalPermissionParams{
				UserID:             userId,
				GlobalPermissionID: "create-projects",
			})

			if err != nil {
				log.Error("project create: auth query error", "err", err)
				return ErrInternal
			}
			if !authorized {
				return ErrNotAuthorizedCreate
			}

			prefix := strings.ToUpper(args[0])
			if len(prefix) != 4 {
				return ErrPrefix4
			}
			title := strings.Join(args[1:], " ")

			err = action.CreateProject(app, query, userId, title, prefix)
			if err != nil {
				log.Errorf("project create: %s", err.Error())
				if err == action.ErrDuplicatePrefix {
					return err
				}
				return ErrInternal
			}

			cmd.Println("project created!")

			err = tx.Commit(ctx)
			if err != nil {
				log.Errorf("project create: %s", err.Error())
				return ErrInternal
			}

			app.Broadcast(projects.RefreshProjectsMsg{})

			return nil
		},
	}

	return createCmd
}
