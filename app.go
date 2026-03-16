package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/log/v2"
	"charm.land/wish/v2"
	"charm.land/wish/v2/activeterm"
	"charm.land/wish/v2/bubbletea"
	"charm.land/wish/v2/logging"
	"github.com/charmbracelet/ssh"
	"github.com/spf13/cobra"
)

// App contains a wish server and the list of running programs.
type App struct {
	*ssh.Server
	progs []*tea.Program
}

// send dispatches a message to all running programs
func (a *App) Broadcast(msg tea.Msg) {
	for _, p := range a.progs {
		go p.Send(msg)
	}
}

func NewApp() *App {
	a := new(App)

	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bubbletea.MiddlewareWithProgramHandler(a.ProgramHandler),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
	}

	a.Server = s
	return a
}

func (a *App) Start() {
	var err error
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server", "host", host, "port", port)
	go func() {
		if err = a.ListenAndServe(); err != nil {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := a.Shutdown(ctx); err != nil {
		log.Error("Could not stop server", "error", err)
	}
}

func cmd(session ssh.Session, app *App, progPtr **tea.Program) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "isshues",
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			model := initialModel()
			model.App = app // each model gets a reference to the global App
			model.id = session.User()

			*progPtr = tea.NewProgram(model, bubbletea.MakeOptions(session)...)
			log.Info("root command called")
			app.progs = append(app.progs, *progPtr)

			return nil
		},
	}
	return cmd
}

func subcmdtest() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "test",
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info("exedcuted test subcommand")
			return nil
		},
	}
	return cmd

}

// function called by wish to create a new tea.Program
func (a *App) ProgramHandler(session ssh.Session) *tea.Program {
	var prog *tea.Program

	rootCmd := cmd(session, a, &prog)
	rootCmd.SetArgs(session.Command())
	rootCmd.SetIn(session)
	rootCmd.SetOut(session)
	rootCmd.SetErr(session.Stderr())
	rootCmd.AddCommand(subcmdtest())
	// rootCmd.CompletionOptions.DisableDefaultCmd = true
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		_ = session.Exit(1)
		return nil
	}

	return prog
}
