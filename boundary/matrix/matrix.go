package matrix

import (
	"context"
	"errors"
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

//go:generate go tool moq -out matrix_mock.go . MatrixBoundary

const (
	HomeserverUrlKey  = "MATRIX_HOMESERVER"
	MatrixUsernameKey = "MATRIX_USERNAME"
	MatrixPasswordKey = "MATRIX_PASSWORD"
	MatrixRoomsKey    = "MATRIX_ROOMS"
)

type MatrixBoundary interface {
	Start()
	loginAndJoin(roomIds []string)
	listen()
	handleMessageEvent(evt *event.Event)
	message(room id.RoomID, content string)
	react(room id.RoomID, evt id.EventID, content string)
	reply(room id.RoomID, evt id.EventID, content string, asHtml bool) id.EventID
	getUserByUsername(tx *gorm.DB, username string) (*entity.User, error)
}

type MatrixBoundaryImpl struct {
	ctx              context.Context
	repo             entity.Repository
	client           *mautrix.Client
	startupTimestamp int64
}

func NewMatrixBoundary(ctx context.Context, repo entity.Repository) *MatrixBoundaryImpl {
	client, err := mautrix.NewClient(ctx.Value(HomeserverUrlKey).(string), "", "")
	if err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("could not create matrix client")
	}
	return &MatrixBoundaryImpl{ctx: ctx, repo: repo, client: client, startupTimestamp: time.Now().UnixMilli()}
}

func (m *MatrixBoundaryImpl) Start() {
	m.loginAndJoin(m.ctx.Value(MatrixRoomsKey).([]string))
	m.listen()
}

func (m *MatrixBoundaryImpl) loginAndJoin(roomIds []string) {
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

func (m *MatrixBoundaryImpl) listen() {
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

func (m *MatrixBoundaryImpl) handleMessageEvent(evt *event.Event) {
	message := evt.Content.AsMessage().Body
	if !strings.HasPrefix(message, ".ordaa") {
		return
	}
	message = strings.TrimSpace(strings.TrimPrefix(message, ".ordaa "))
	log.Ctx(m.ctx).Debug().Msgf("received message: %s", message)

	err := m.repo.Transaction(func(tx *gorm.DB) error {
		commands := strings.Split(message, " ")
		if len(commands) < 0 {
			return handleUnrecognizedCommand(m.ctx, m, m.repo, tx, evt, message)
		}
		command := commands[0]

		handler := handlers[command]
		if handler == nil {
			return handleUnrecognizedCommand(m.ctx, m, m.repo, tx, evt, message)
		}

		return handler(m.ctx, m, m.repo, tx, evt, message)
	})
	if err != nil {
		m.reply(evt.RoomID, evt.ID, err.Error(), false)
	}
}

func (m *MatrixBoundaryImpl) message(room id.RoomID, content string) {
	_, err := m.client.SendNotice(m.ctx, room, content)
	if err != nil {
		log.Ctx(m.ctx).Error().Err(err).Msgf("could not send message '%s'", content)
	}
}

func (m *MatrixBoundaryImpl) react(room id.RoomID, evt id.EventID, content string) {
	_, err := m.client.SendReaction(m.ctx, room, evt, content)
	if err != nil {
		log.Ctx(m.ctx).Error().Err(err).Msg("could not react to event")
	}
}

func (m *MatrixBoundaryImpl) reply(room id.RoomID, evt id.EventID, content string, asHtml bool) id.EventID {
	contentJSON := map[string]interface{}{
		"m.relates_to": map[string]interface{}{
			"m.in_reply_to": map[string]interface{}{
				"event_id": evt,
			},
		},
		"msgtype": "m.text",
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

func (m *MatrixBoundaryImpl) getUserByUsername(tx *gorm.DB, username string) (*entity.User, error) {
	matrixUser, err := m.repo.GetMatrixUserByUsername(tx, username)
	if err != nil {
		msg := fmt.Sprintf("could not get matrix user for username '%s'", username)
		log.Ctx(m.ctx).Warn().Err(err).Msg(msg)
		return nil, errors.New(msg)
	}

	user, err := m.repo.GetUser(tx, matrixUser.UserUuid)
	if err != nil {
		msg := fmt.Sprintf("could not get user for matrix user '%s'", matrixUser.Username)
		log.Ctx(m.ctx).Warn().Err(err).Msg(msg)
		return nil, errors.New(msg)
	}

	return user, nil
}
