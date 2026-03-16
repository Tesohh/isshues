package app

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/log/v2"
	"charm.land/wish/v2/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/spf13/cobra"
)

// a isshues command can modify the progPtr to initialize and give a tea.Program to the program handler.
//
// if no tea.Program needs to be created (eg. cli only actions) leave it as nil
type IsshuesCommand func(session ssh.Session, app *App, progPtr **tea.Program) *cobra.Command

func cmd(session ssh.Session, app *App, progPtr **tea.Program) *cobra.Command {
	return &cobra.Command{
		Use:  "isshues",
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			model := initialModel()
			model.app = app // each model gets a reference to the global App
			model.id = session.User()

			*progPtr = tea.NewProgram(model, bubbletea.MakeOptions(session)...)
			log.Info("root command called")
			app.progs = append(app.progs, *progPtr)

			return nil
		},
	}
}
