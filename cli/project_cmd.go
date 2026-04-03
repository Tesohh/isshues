package cli

import (
	"errors"

	tea "charm.land/bubbletea/v2"
	"github.com/Tesohh/isshues/app"
	"github.com/charmbracelet/ssh"
	"github.com/spf13/cobra"
)

var (
	NotAuthorizedCreateErr = errors.New("you are not allowed to do this. your account is missing the create-projects global permission")
	InternalErr            = errors.New("internal error. please contact your admin.")
	Prefix4Err             = errors.New("the prefix must be 4 characters long")
	DuplicatePrefixErr     = errors.New("this prefix is already taken")
)

func projectCmd(session ssh.Session, app *app.App, program **tea.Program) *cobra.Command {
	projectCmd := &cobra.Command{
		Use: "project",
	}

	projectCmd.AddCommand(projectCreateCmd(session, app, program))

	return projectCmd
}
