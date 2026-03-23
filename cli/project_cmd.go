package cli

import (
	"context"
	"errors"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/log/v2"
	"github.com/Tesohh/isshues/app"
	"github.com/Tesohh/isshues/config"
	db "github.com/Tesohh/isshues/db/generated"
	"github.com/Tesohh/isshues/model/projects"
	"github.com/charmbracelet/ssh"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/spf13/cobra"
)

var (
	NotAuthorizedCreateErr = errors.New("you are not allowed to do this. your account is missing the create-projects global permission")
	InternalDbErr          = errors.New("internal db error. please contact your admin.")
	InternalConfigErr      = errors.New("internal config error. please contact your admin.")
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
				return InternalDbErr
			}
			if !authorized {
				return NotAuthorizedCreateErr
			}

			prefix := strings.ToUpper(args[0])
			if len(prefix) != 4 {
				return Prefix4Err
			}
			title := strings.Join(args[1:], " ")

			// create the project
			projectId, err := app.GetDB().InsertProject(ctx, db.InsertProjectParams{
				Title:  title,
				Prefix: prefix,
			})
			var pgerr *pgconn.PgError
			if errors.As(err, &pgerr) && pgerr.Code == "23505" {
				return DuplicatePrefixErr
			} else if err != nil {
				log.Error("project new: insertion error", "err", err)
				return InternalDbErr
			}

			// add default groups
			// TODO: get these from the server config

			var defaultGroups []config.DefaultGroup
			err = app.Viper.UnmarshalKey("default_groups", &defaultGroups)
			if err != nil {
				log.Error("project new: default_groups config error", "err", err)
				return InternalConfigErr
			}

			for _, group := range defaultGroups {
				params := db.InsertGroupParams{
					Name:        pgtype.Text{String: group.Name, Valid: true},
					Color:       pgtype.Text{String: group.Color, Valid: group.Color != ""},
					Mentionable: group.Mentionable,
					ProjectID:   projectId,
				}
				groupId, err := app.GetDB().InsertGroup(ctx, params)
				if err != nil {
					log.Error("project new: group insertion error", "err", err)
					return InternalDbErr
				}

				for _, permission := range group.Permissions {
					err := app.GetDB().GrantPermissionToGroup(ctx, db.GrantPermissionToGroupParams{
						GroupID:             groupId,
						ProjectPermissionID: permission,
					})
					if err != nil {
						log.Error("project new: group grant permission error", "err", err, "permission", permission)
						return InternalDbErr
					}
				}

				// add creator to the admins group
				if group.AddCreator {
					err = app.GetDB().AddUserToGroup(ctx, db.AddUserToGroupParams{
						GroupID: groupId,
						UserID:  userId,
					})
					if err != nil {
						log.Error("project new: cannot add creator to group", "err", err)
						return InternalDbErr
					}
				}
			}

			cmd.Println("project created!")
			app.Broadcast(projects.RefreshProjectsMsg{})

			return nil
		},
	}

	projectCmd.AddCommand(newCmd)

	return projectCmd
}
