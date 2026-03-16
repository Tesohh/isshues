package app

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/log/v2"
	"github.com/charmbracelet/ssh"
	"github.com/spf13/cobra"
)

func subcmdtest(_ ssh.Session, _ *App, _ **tea.Program) *cobra.Command {
	return &cobra.Command{
		Use:  "test",
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Println("Hello from the server!")
			log.Info("executed test subcommand")
			return nil
		},
	}
}
