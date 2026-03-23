package app

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
	db "github.com/Tesohh/isshues/db/generated"
	"github.com/charmbracelet/ssh"
	"github.com/jackc/pgx/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// App contains a wish server and the list of running programs.
type App struct {
	*ssh.Server
	host string
	port string

	Viper *viper.Viper
	DB    *db.Queries

	SessionIdToUserIds map[string]int64

	Progs []*tea.Program
}

func (a *App) GetDB() *db.Queries {
	return a.DB
}

// send dispatches a message to all running programs
// TODO: consider making an interface for this to be used in models in another packages.
func (a *App) Broadcast(msg tea.Msg) {
	for _, p := range a.Progs {
		go p.Send(msg)
	}
}

type isshuesCmd func(session ssh.Session, app *App, progPtr **tea.Program) *cobra.Command

func NewApp(dbConn *pgx.Conn, viper *viper.Viper, rootCmd isshuesCmd) *App {
	a := new(App)
	a.SessionIdToUserIds = make(map[string]int64)

	a.Viper = viper
	a.host = viper.GetString("ssh.host")
	a.port = viper.GetString("ssh.port")
	a.DB = db.New(dbConn)

	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(a.host, a.port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithPublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
			return key.Type() == "ssh-ed25519"
		}),
		wish.WithMiddleware(
			bubbletea.MiddlewareWithProgramHandler(a.MakeProgramHandler(rootCmd)),
			a.AuthMiddleware,
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
	log.Info("Starting SSH server", "host", a.host, "port", a.port)
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

// function called by wish to create a new tea.Program
func (a *App) MakeProgramHandler(rootCmd isshuesCmd) func(session ssh.Session) *tea.Program {
	return func(session ssh.Session) *tea.Program {
		var prog *tea.Program

		session.PublicKey()

		rootCmd := rootCmd(session, a, &prog)
		rootCmd.SetArgs(session.Command())
		rootCmd.SetIn(session)
		rootCmd.SetOut(session)
		rootCmd.SetErr(session.Stderr())
		rootCmd.CompletionOptions.DisableDefaultCmd = true
		if err := rootCmd.Execute(); err != nil {
			log.Error(err)
			_ = session.Exit(1)
			return nil
		}

		return prog
	}
}
