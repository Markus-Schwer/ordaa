package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

type MatrixBot struct {
	ctx     context.Context
	client  *mautrix.Client
	startup int64
	parser  *MessageParser
	runner  *ActionRunner
}

// assumes home server to be 'matrix.org'
func NewMatrixBot(ctx context.Context, parser *MessageParser, runner *ActionRunner) *MatrixBot {
	client, err := mautrix.NewClient(ctx.Value(HomeServerURLKey).(string), "", "")
	if err != nil {
		log.Fatal().Err(err).Msg("could not create matrix client")
	}
	return &MatrixBot{
		ctx:     ctx,
		client:  client,
		startup: time.Now().UnixMilli(),
		parser:  parser,
		runner:  runner,
	}
}

func (bot *MatrixBot) LoginAndJoin(roomIds []string) {
	log.Debug().Msg("Logging in")
	_, err := bot.client.Login(&mautrix.ReqLogin{
		Type:               mautrix.AuthTypePassword,
		Identifier:         mautrix.UserIdentifier{Type: mautrix.IdentifierTypeUser, User: bot.ctx.Value(UserKey).(string)},
		Password:           bot.ctx.Value(PasswordKey).(string),
		StoreCredentials:   true,
		StoreHomeserverURL: true,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("could not do login")
	}
	err = bot.client.SetDisplayName("Chicken Masalla legende Wollmilchsau [BOT]")
	if err != nil {
		log.Err(err)
	}
	for _, roomId := range roomIds {
		_, err = bot.client.JoinRoomByID(id.RoomID(roomId))
		if err != nil {
			log.Fatal().Err(err).Msg("could not join room")
		}
		log.Debug().Msgf("joined room %s", roomId)
	}
}

func (bot *MatrixBot) Listen() {
	log.Debug().Msg("listening to messages")
	syncer := bot.client.Syncer.(*mautrix.DefaultSyncer)
	syncer.OnEventType(event.EventMessage, func(src mautrix.EventSource, evt *event.Event) {
		bot.handleMessageEvent(evt)
	})
	syncer.OnEventType(event.EventReaction, func(src mautrix.EventSource, evt *event.Event) {
		bot.handleReactionEvent(evt)
	})
	err := bot.client.SyncWithContext(bot.ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("client had a problem when syncing")
	}
}

func (bot *MatrixBot) handleMessageEvent(evt *event.Event) {
	if evt.Timestamp < bot.startup || evt.Sender == bot.client.UserID {
		return
	}
	if !strings.HasPrefix(evt.Content.AsMessage().Body, bot.parser.trigger) {
		return
	}
	msgAction, err := bot.parser.convertToAction(evt.Content.AsMessage().Body)
	if err != nil {
		log.Err(err)
		bot.reply(evt.RoomID, evt.ID, err.Error())
		return
	}
	err = bot.runner.runAction(evt.Sender.String(), msgAction)
	if err != nil {
		log.Err(err)
		bot.reply(evt.RoomID, evt.ID, err.Error())
		return
	}
	bot.react(evt.RoomID, evt.ID, "✅")
	log.Debug().Msgf("handled actions for message '%s'", evt.Content.AsMessage().Body)
}

func (bot *MatrixBot) handleReactionEvent(evt *event.Event) {
	if evt.Timestamp < bot.startup || evt.Sender == bot.client.UserID {
		return
	}
	// get message to which the reaction relates
	msgEvt, err := bot.client.GetEvent(evt.RoomID, evt.Content.AsReaction().RelatesTo.EventID)
	if err != nil {
		log.Error().Err(err).Msg("could not get the message the reaction relates to")
		return
	}
	// parsing has to be triggered manually when `GetEvent` is used
	err = msgEvt.Content.ParseRaw(event.EventMessage)
	if err != nil {
		log.Err(err)
		return
	}
	// get the initial action to use it as base
	action, err := bot.parser.convertToAction(msgEvt.Content.AsMessage().Body)
	switch evt.Content.AsReaction().RelatesTo.Key {
	// run the same action again for the user who reacted
	case "📈":
		err = bot.runner.runAction(evt.Sender.String(), action)
		if err != nil {
			log.Err(err)
			return
		}
		bot.reply(evt.RoomID, msgEvt.ID, fmt.Sprintf("%s: ✅", evt.Sender.String()))
		return
	// remove the item, only possible on own actions
	case "📉":
		if msgEvt.Sender != evt.Sender {
			log.Warn().Msgf("user '%s' tried to troll user '%s'", evt.Sender.String(), msgEvt.Sender.String())
			bot.message(evt.RoomID, fmt.Sprintf("hilarious idea @%s:, but no", evt.Sender.String()))
			return
		}
		action.verb = Remove
		err = bot.runner.runAction(evt.Sender.String(), action)
		if err != nil {
			log.Err(err)
			return
		}
		bot.reply(evt.RoomID, msgEvt.ID, fmt.Sprintf("%s: ✅", evt.Sender.String()))
		return
	}
}

func (bot *MatrixBot) message(room id.RoomID, content string) {
	_, err := bot.client.SendNotice(room, content)
	if err != nil {
		log.Error().Err(err).Msgf("could not send message '%s'", content)
	}
}

func (bot *MatrixBot) react(room id.RoomID, evt id.EventID, content string) {
	_, err := bot.client.SendReaction(room, evt, content)
	if err != nil {
		log.Error().Err(err).Msg("could not react to event")
	}
}

func (bot *MatrixBot) reply(room id.RoomID, evt id.EventID, content string) {
	_, err := bot.client.SendMessageEvent(room, event.EventMessage, map[string]interface{}{
		"m.relates_to": map[string]interface{}{
			"m.in_reply_to": map[string]interface{}{
				"event_id": evt,
			},
		},
		// notice is a message from a bot, it avoids feedback loops
		"msgtype": "m.notice",
		"body":    content,
	})
	if err != nil {
		log.Error().Err(err).Msgf("could not respond '%s' to event", content)
	}
}
