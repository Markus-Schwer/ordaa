package tui

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	"github.com/rs/zerolog/log"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/entity"
)

var users = map[string]string{
	"jgero": "sk-ssh-ed25519@openssh.com AAAAGnNrLXNzaC1lZDI1NTE5QG9wZW5zc2guY29tAAAAIKgKBq4N0toosQ6nV/IQTRx/8OudkB7DwnrIDX0HrUw7AAAABHNzaDo=",
}

func NewSshTuiServer(ctx context.Context, repo entity.Repository) *SshTuiServer {
	return &SshTuiServer{
		host: "localhost",
		port: "23234",
		ctx:  ctx,
		repo: repo,
	}
}

type SshTuiServer struct {
	host string
	port string
	ctx  context.Context
	repo entity.Repository
}

func (serv *SshTuiServer) Start() error {
	zlog := log.Ctx(serv.ctx).With().Str("component", "ssh-server").Logger()
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(serv.host, serv.port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithPublicKeyAuth(func(_ ssh.Context, key ssh.PublicKey) bool {
			for _, pubkey := range users {
				parsed, _, _, _, _ := ssh.ParseAuthorizedKey(
					[]byte(pubkey),
				)
				if ssh.KeysEqual(key, parsed) {
					return true
				}
			}
			return false
		}),
		wish.WithMiddleware(
			bubbletea.Middleware(serv.teaHandler),
			activeterm.Middleware(), // Bubble Tea apps usually require a PTY.
			logging.MiddlewareWithLogger(&zlog),
		),
	)
	if err != nil {
		log.Ctx(serv.ctx).Error().Err(err).Msg("could not start server")
		return err
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Ctx(serv.ctx).Info().Msgf("starting ssh server at '%v:%v'", serv.host, serv.port)
	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Ctx(serv.ctx).Error().Err(err).Msg("could not start server")
			done <- nil
		}
	}()

	<-done
	log.Ctx(serv.ctx).Info().Msg("stopping sh server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Ctx(serv.ctx).Error().Err(err).Msg("could not stop server")
		return err
	}
	return nil
}

func (serv *SshTuiServer) teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, _ := s.Pty()

	renderer := bubbletea.MakeRenderer(s)
	m := NewLayoutModel(serv.ctx, renderer, pty, serv.repo)

	return m, []tea.ProgramOption{tea.WithAltScreen()}
}
