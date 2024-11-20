package tui

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
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
	"gorm.io/gorm"
)

const UserContextKey = "user-struct-context-key"

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
		wish.WithPublicKeyAuth(func(_ ssh.Context, _ ssh.PublicKey) bool { return true }),
		wish.WithMiddleware(
			bubbletea.Middleware(serv.teaHandler),
			activeterm.Middleware(), // Bubble Tea apps usually require a PTY.
			logging.MiddlewareWithLogger(&zlog),
			func(next ssh.Handler) ssh.Handler {
				return func(s ssh.Session) {
					user, err := serv.getUserForKey(s)
					if err != nil {
						log.Ctx(serv.ctx).Error().Err(err).Msg("could not get user for key")
						return
					}
					s.Context().SetValue(UserContextKey, user)
					next(s)
				}
			},
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

func (serv *SshTuiServer) getUserForKey(sess ssh.Session) (user *entity.User, err error) {
	err = serv.repo.Transaction(func(tx *gorm.DB) error {
		var innerErr error
		if sess.PublicKey() == nil {
			return fmt.Errorf("could not get public key from session: %w", err)
		}
		pubKey := fmt.Sprintf("%s %s",sess.PublicKey().Type(), base64.StdEncoding.EncodeToString(sess.PublicKey().Marshal()))
		sshUser, innerErr := serv.repo.GetSshUserByPublicKey(tx, pubKey)
		if innerErr != nil {
			return innerErr
		}

		user, innerErr = serv.repo.GetUser(tx, sshUser.UserUuid)
		if innerErr != nil {
			return innerErr
		}

		return nil
	})
	if err != nil {
		return
	}
	return
}
