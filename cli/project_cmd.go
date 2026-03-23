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

var (
	NotAuthorizedCreateErr = errors.New("you are not allowed to do this. your account is missing the create-projects global permission")
	InternalErr            = errors.New("internal error. please contact your admin.")
	Prefix4Err             = errors.New("the prefix must be 4 characters long")
	DuplicatePrefixErr     = errors.New("this prefix is already taken")
)

func projectCmd(session ssh.Session, app *app.App, _ **tea.Program) *cobra.Command {
	projectCmd := &cobra.Command{
		Use: "project",
	}

	newCmd := &cobra.Command{
		Use:  "new [prefix] [title]",
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			userId, ok := app.SessionIdToUserIds[session.Context().SessionID()]
			if !ok {
				return errors.New("your userid was not found in the session map. might be an auth issue.")
			}

			ctx := context.Background()
			authorized, err := app.GetDB().UserHasGlobalPermission(ctx, db.UserHasGlobalPermissionParams{
				UserID:             userId,
				GlobalPermissionID: "create-projects",
			})

			if err != nil {
				log.Error("project new: auth query error", "err", err)
				return InternalErr
			}
			if !authorized {
				return NotAuthorizedCreateErr
			}

			prefix := strings.ToUpper(args[0])
			if len(prefix) != 4 {
				return Prefix4Err
			}
			title := strings.Join(args[1:], " ")

			err = action.CreateProject(app, userId, title, prefix)
			if err != nil {
				log.Errorf("project new: %s", err.Error())
				if err == action.DuplicatePrefixErr {
					return err
				}
				return InternalErr
			}

			cmd.Println("project created!")
			app.Broadcast(projects.RefreshProjectsMsg{})

			return nil
		},
	}

	projectCmd.AddCommand(newCmd)

	return projectCmd
}
