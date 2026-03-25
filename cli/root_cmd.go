package cli

import (
	"errors"

	tea "charm.land/bubbletea/v2"
	"charm.land/log/v2"
	"charm.land/wish/v2/bubbletea"
	"github.com/Tesohh/isshues/app"
	"github.com/Tesohh/isshues/model/root"
	"github.com/charmbracelet/ssh"
	"github.com/spf13/cobra"
)

func RootCmd(session ssh.Session, app *app.App, progPtr **tea.Program) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "isshues",
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			userId, ok := app.SessionIdToUserIds[session.Context().SessionID()]
			if !ok {
				return errors.New("your userid was not found in the session map. might be an auth issue.")
			}

			model := root.New(app, userId)
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
