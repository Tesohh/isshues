package cli

import (
	"context"
	"errors"
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/log/v2"
	"charm.land/wish/v2/bubbletea"
	"github.com/Tesohh/isshues/app"
	"github.com/Tesohh/isshues/model/root"
	"github.com/charmbracelet/ssh"
	tint "github.com/lrstanley/bubbletint/v2"
	"github.com/spf13/cobra"
)

var (
	ThemeNotFoundErr = errors.New("theme not found")
)

func RootCmd(session ssh.Session, app *app.App, progPtr **tea.Program) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "isshues",
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			userId, ok := app.SessionIdToUserIds[session.Context().SessionID()]
			if !ok {
				return errors.New("your userid was not found in the session map. might be an auth issue.")
			}

			// // TODO: don't hardcode this
			// theme, _ := tint.GetTint("gruvbox_dark")

			settings, err := app.DB.GetUserSettings(ctx, userId)
			if err != nil {
				log.Error("settings query error", "err", err, "userId", userId)
				return InternalErr
			}

			theme, ok := tint.GetTint(settings.Theme)
			if !ok {
				return fmt.Errorf("%w: %s. go to https://lrstanley.github.io/bubbletint/ to find a list of supported themes", ThemeNotFoundErr, settings.Theme)
			}

			model := root.New(app, userId, theme)
			model.App = app // each model gets a reference to the global App

			*progPtr = tea.NewProgram(model, bubbletea.MakeOptions(session)...)
			log.Info("root command called")
			app.Progs = append(app.Progs, *progPtr)

			return nil
		},
	}

	cmd.AddCommand(projectCmd(session, app, progPtr))
	return cmd
}
