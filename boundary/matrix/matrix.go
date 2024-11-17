package matrix

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/entity"
	"gorm.io/gorm"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

const (
	HomeserverUrlKey  = "MATRIX_HOMESERVER"
	MatrixUsernameKey = "MATRIX_USERNAME"
	MatrixPasswordKey = "MATRIX_PASSWORD"
	MatrixRoomsKey    = "MATRIX_ROOMS"
)

type MatrixBoundary struct {
	ctx              context.Context
	repo             entity.Repository
	client           *mautrix.Client
	startupTimestamp int64
}

func NewMatrixBoundary(ctx context.Context, repo entity.Repository) *MatrixBoundary {
	client, err := mautrix.NewClient(ctx.Value(HomeserverUrlKey).(string), "", "")
	if err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("could not create matrix client")
	}
	return &MatrixBoundary{ctx: ctx, repo: repo, client: client, startupTimestamp: time.Now().UnixMilli()}
}

func (m *MatrixBoundary) Start() {
	m.loginAndJoin(m.ctx.Value(MatrixRoomsKey).([]string))
	m.listen()
}

func (m *MatrixBoundary) loginAndJoin(roomIds []string) {
	log.Ctx(m.ctx).Debug().Msg("Logging in to matrix homeserver")
	_, err := m.client.Login(m.ctx, &mautrix.ReqLogin{
		Type:               mautrix.AuthTypePassword,
		Identifier:         mautrix.UserIdentifier{Type: mautrix.IdentifierTypeUser, User: m.ctx.Value(MatrixUsernameKey).(string)},
		Password:           m.ctx.Value(MatrixPasswordKey).(string),
		StoreCredentials:   true,
		StoreHomeserverURL: true,
	})
	if err != nil {
		log.Ctx(m.ctx).Fatal().Err(err).Msg("could not login")
	}

	err = m.client.SetDisplayName(m.ctx, "Chicken Masalla legende Wollmilchsau [m]")
	if err != nil {
		log.Ctx(m.ctx).Err(err)
	}

	for _, roomId := range roomIds {
		_, err = m.client.JoinRoomByID(m.ctx, id.RoomID(roomId))
		if err != nil {
			log.Ctx(m.ctx).Fatal().Err(err).Msg("could not join room")
		}
		log.Ctx(m.ctx).Debug().Msgf("joined room %s", roomId)
	}
}

func (m *MatrixBoundary) listen() {
	log.Ctx(m.ctx).Debug().Msg("listening to matrix messages")
	syncer := m.client.Syncer.(*mautrix.DefaultSyncer)
	syncer.OnEventType(event.EventMessage, func(ctx context.Context, evt *event.Event) {
		if evt.Timestamp < m.startupTimestamp || evt.Sender == m.client.UserID {
			return
		}
		m.handleMessageEvent(evt)
	})

	err := m.client.SyncWithContext(m.ctx)
	if err != nil {
		log.Ctx(m.ctx).Fatal().Err(err).Msg("client had a problem when syncing")
	}
}

func (m *MatrixBoundary) handleMessageEvent(evt *event.Event) {
	message := evt.Content.AsMessage().Body
	if !strings.HasPrefix(message, ".ordaa") {
		return
	}
	message = strings.TrimSpace(strings.TrimPrefix(message, ".ordaa "))
	log.Ctx(m.ctx).Debug().Msgf("received message: %s", message)

	err := m.repo.Transaction(func(tx *gorm.DB) error {
		commands := strings.Split(message, " ")
		if len(commands) < 0 {
			return handleUnrecognizedCommand(m, tx, evt, message)
		}
		command := commands[0]

		handler := handlers[command]
		if handler == nil {
			return handleUnrecognizedCommand(m, tx, evt, message)
		}

		return handler(m, tx, evt, message)
	})
	if err != nil {
		log.Ctx(m.ctx).Error().Err(err).Msgf("error occured while handling matrix message: %s", message)
		m.reply(evt.RoomID, evt.ID, fmt.Sprintf("error occured while handling matrix: %s", err), false)
	}
}

func (m *MatrixBoundary) message(room id.RoomID, content string) {
	_, err := m.client.SendNotice(m.ctx, room, content)
	if err != nil {
		log.Ctx(m.ctx).Error().Err(err).Msgf("could not send message '%s'", content)
	}
}

func (m *MatrixBoundary) react(room id.RoomID, evt id.EventID, content string) {
	_, err := m.client.SendReaction(m.ctx, room, evt, content)
	if err != nil {
		log.Ctx(m.ctx).Error().Err(err).Msg("could not react to event")
	}
}

func (m *MatrixBoundary) reply(room id.RoomID, evt id.EventID, content string, asHtml bool) id.EventID {
	contentJSON := map[string]interface{}{
		"m.relates_to": map[string]interface{}{
			"m.in_reply_to": map[string]interface{}{
				"event_id": evt,
			},
		},
		// notice is a message from a bot, it avoids feedback loops
		"msgtype": "m.notice",
		"body":    content,
	}
	if asHtml {
		contentJSON["format"] = "org.matrix.custom.html"
		contentJSON["formatted_body"] = strings.TrimSuffix(content, "\n")
	}
	ev, err := m.client.SendMessageEvent(m.ctx, room, event.EventMessage, contentJSON)
	if err != nil {
		log.Ctx(m.ctx).Error().Err(err).Msgf("could not respond to event '%s'", content)
	}
	return ev.EventID
}

func (m *MatrixBoundary) getUserByUsername(tx *gorm.DB, username string) (*entity.User, error) {
	matrixUser, err := m.repo.GetMatrixUserByUsername(tx, username)
	if err != nil {
		return nil, fmt.Errorf("could not get matrix user of sender '%s': %w", username, err)
	}

	user, err := m.repo.GetUser(tx, matrixUser.UserUuid)
	if err != nil {
		return nil, fmt.Errorf("could not get user of sender '%s' for message: %w", username, err)
	}

	return user, nil
}
