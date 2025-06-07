package matrix

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"

	"github.com/Markus-Schwer/ordaa/internal/boundary/matrix/handler"
	"github.com/Markus-Schwer/ordaa/internal/config"
)

var ErrGettingDefaultSyncer = errors.New("getting DefaultSyncer")

type CommandHandler interface {
	Matches(ctx context.Context, evt *event.Event) bool
	Handle(ctx context.Context, evt *event.Event) *handler.CommandResponse
}

type Boundary struct {
	cfg              *config.MatrixConfig
	client           *mautrix.Client
	startupTimestamp int64
	handlers         []CommandHandler
}

func NewMatrixBoundary(
	ctx context.Context,
	cfg *config.MatrixConfig,
	userService handler.UserService,
	orderService handler.OrderService,
) (*Boundary, error) {
	client, err := mautrix.NewClient(cfg.HomeserverURL, "", "")
	if err != nil {
		return nil, fmt.Errorf("creating matrix client: %w", err)
	}

	return &Boundary{
		cfg:              cfg,
		client:           client,
		startupTimestamp: time.Now().UnixMilli(),
		handlers: []CommandHandler{
			&handler.HelpHandler{},
			&handler.RegisterHandler{UserService: userService},
			&handler.StartHandler{UserService: userService, OrderService: orderService},
			&handler.AddHandler{UserService: userService, OrderService: orderService},
			&handler.StateTransitionHandler{UserService: userService, OrderService: orderService},
			&handler.UnrecognizedCommandHandler{}, // must be last handler in list, because it always matches
		},
	}, nil
}

func (m *Boundary) Start(ctx context.Context) error {
	if err := m.loginAndJoin(ctx, m.cfg.Rooms); err != nil {
		return err
	}

	return m.listen(ctx)
}

func (m *Boundary) Stop() error {
	log.Info().Msg("shutting down matrix boundary")

	m.client.StopSync()

	return nil
}

func (m *Boundary) loginAndJoin(ctx context.Context, roomIDs []string) error {
	log.Ctx(ctx).Debug().Msg("Logging in to matrix homeserver")

	req := &mautrix.ReqLogin{
		Type:               mautrix.AuthTypePassword,
		Identifier:         mautrix.UserIdentifier{Type: mautrix.IdentifierTypeUser, User: m.cfg.Username},
		Password:           m.cfg.Password,
		StoreCredentials:   true,
		StoreHomeserverURL: true,
	}

	if _, err := m.client.Login(ctx, req); err != nil {
		return fmt.Errorf("logging in to matrix homeserver: %w", err)
	}

	if err := m.client.SetDisplayName(ctx, m.cfg.DisplayName); err != nil {
		return fmt.Errorf("setting display name: %w", err)
	}

	for _, roomID := range roomIDs {
		if _, err := m.client.JoinRoomByID(ctx, id.RoomID(roomID)); err != nil {
			return fmt.Errorf("joining room: %w", err)
		}

		log.Ctx(ctx).Debug().Msgf("joined room %s", roomID)
	}

	return nil
}

func (m *Boundary) listen(ctx context.Context) error {
	log.Ctx(ctx).Debug().Msg("listening to matrix events")

	syncer, ok := m.client.Syncer.(*mautrix.DefaultSyncer)
	if !ok {
		return ErrGettingDefaultSyncer
	}

	syncer.OnEventType(event.EventMessage, func(ctx context.Context, evt *event.Event) {
		if evt.Timestamp < m.startupTimestamp || evt.Sender == m.client.UserID {
			return
		}

		m.handleMessageEvent(ctx, evt)
	})

	if err := m.client.SyncWithContext(ctx); err != nil {
		return fmt.Errorf("listening to matrix events: %w", err)
	}

	return nil
}

func (m *Boundary) handleMessageEvent(ctx context.Context, evt *event.Event) {
	msg := evt.Content.AsMessage().Body
	if !strings.HasPrefix(msg, handler.MatrixCommandPrefix) {
		return
	}

	log.Ctx(ctx).Debug().Msgf("received message: %s", msg)

	for _, handler := range m.handlers {
		if !handler.Matches(ctx, evt) {
			continue
		}

		resp := handler.Handle(ctx, evt)
		if resp == nil {
			log.Ctx(ctx).Warn().Msgf("command handler didn't return a response for message: %s", msg)
			break
		}

		if err := m.reply(ctx, evt.RoomID, evt.ID, resp.Msg, resp.AsHTML); err != nil {
			log.Ctx(ctx).Warn().Err(err).Msg("handling command")
		}

		break
	}
}

// func (m *MatrixBoundary) message(ctx context.Context, room id.RoomID, content string) error {
//	if _, err := m.client.SendNotice(ctx, room, content); err != nil {
//		return fmt.Errorf("sending message: %w", err)
//	}
//
//	return nil
//}

// func (m *MatrixBoundary) react(ctx context.Context, room id.RoomID, evt id.EventID, content string) error {
//	if _, err := m.client.SendReaction(ctx, room, evt, content); err != nil {
//		return fmt.Errorf("sending reaction to event: %w", err)
//	}
//
//	return nil
//}

func (m *Boundary) reply(ctx context.Context, room id.RoomID, evt id.EventID, content string, asHTML bool) error {
	contentJSON := map[string]any{
		"m.relates_to": map[string]any{
			"m.in_reply_to": map[string]any{
				"event_id": evt,
			},
		},
		"msgtype": "m.text",
		"body":    content,
	}

	if asHTML {
		contentJSON["format"] = "org.matrix.custom.html"
		contentJSON["formatted_body"] = strings.TrimSuffix(content, "\n")
	}

	if _, err := m.client.SendMessageEvent(ctx, room, event.EventMessage, contentJSON); err != nil {
		return fmt.Errorf("responding to event: %w", err)
	}

	return nil
}
