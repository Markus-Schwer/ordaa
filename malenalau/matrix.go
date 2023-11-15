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

type orderMessage struct {
	userId   string
	itemId   string
	provider string
}

type MatrixBot struct {
	ctx          context.Context
	client       *mautrix.Client
	startup      int64
	parser       *MessageParser
	runner       *ActionRunner
	trackedEvent map[string]orderMessage
}

// assumes home server to be 'matrix.org'
func NewMatrixBot(ctx context.Context, parser *MessageParser, runner *ActionRunner) *MatrixBot {
	client, err := mautrix.NewClient(ctx.Value(HomeServerURLKey).(string), "", "")
	if err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("could not create matrix client")
	}
	return &MatrixBot{
		ctx:          ctx,
		client:       client,
		startup:      time.Now().UnixMilli(),
		parser:       parser,
		runner:       runner,
		trackedEvent: make(map[string]orderMessage),
	}
}

func (bot *MatrixBot) LoginAndJoin(roomIds []string) {
	log.Ctx(bot.ctx).Debug().Msg("Logging in")
	_, err := bot.client.Login(&mautrix.ReqLogin{
		Type:               mautrix.AuthTypePassword,
		Identifier:         mautrix.UserIdentifier{Type: mautrix.IdentifierTypeUser, User: bot.ctx.Value(UserKey).(string)},
		Password:           bot.ctx.Value(PasswordKey).(string),
		StoreCredentials:   true,
		StoreHomeserverURL: true,
	})
	if err != nil {
		log.Ctx(bot.ctx).Fatal().Err(err).Msg("could not do login")
	}
	err = bot.client.SetDisplayName("Chicken Masalla legende Wollmilchsau [BOT]")
	if err != nil {
		log.Ctx(bot.ctx).Err(err)
	}
	for _, roomId := range roomIds {
		_, err = bot.client.JoinRoomByID(id.RoomID(roomId))
		if err != nil {
			log.Ctx(bot.ctx).Fatal().Err(err).Msg("could not join room")
		}
		log.Ctx(bot.ctx).Debug().Msgf("joined room %s", roomId)
	}
}

func (bot *MatrixBot) Listen() {
	log.Ctx(bot.ctx).Debug().Msg("listening to messages")
	syncer := bot.client.Syncer.(*mautrix.DefaultSyncer)
	syncer.OnEventType(event.EventMessage, func(src mautrix.EventSource, evt *event.Event) {
		if evt.Timestamp < bot.startup || evt.Sender == bot.client.UserID {
			return
		}
		bot.handleMessageEvent(evt)
	})
	syncer.OnEventType(event.EventReaction, func(src mautrix.EventSource, evt *event.Event) {
		if evt.Timestamp < bot.startup || evt.Sender == bot.client.UserID {
			return
		}
		bot.handleReactionEvent(evt)
	})
	syncer.OnEventType(event.EventRedaction, func(src mautrix.EventSource, evt *event.Event) {
		if evt.Timestamp < bot.startup || evt.Sender == bot.client.UserID {
			return
		}
		bot.handleRedactionEvent(evt)
	})
	err := bot.client.SyncWithContext(bot.ctx)
	if err != nil {
		log.Ctx(bot.ctx).Fatal().Err(err).Msg("client had a problem when syncing")
	}
}

func (bot *MatrixBot) handleRedactionEvent(evt *event.Event) {
	trackedEv := bot.trackedEvent[evt.Redacts.String()]
	if &trackedEv == nil {
		return
	}
	_, err := bot.runner.runAction(trackedEv.userId, &ParsedAction{trackedEv.provider, Remove, trackedEv.itemId})
	if err != nil {
		log.Ctx(bot.ctx).Err(err)
	}
	bot.reply(evt.RoomID, evt.ID, fmt.Sprintf("%s: ✅ deleted %s", trackedEv.userId, trackedEv.itemId), false)
}

func (bot *MatrixBot) handleMessageEvent(evt *event.Event) {
	if !strings.HasPrefix(evt.Content.AsMessage().Body, bot.parser.trigger) {
		return
	}
	msgAction, err := bot.parser.convertToAction(evt.Content.AsMessage().Body)
	if err != nil {
		log.Ctx(bot.ctx).Err(err)
		bot.reply(evt.RoomID, evt.ID, err.Error(), false)
		return
	}
	var message []string
	message, err = bot.runner.runAction(evt.Sender.String(), msgAction)
	if msgAction.verb == Add {
		bot.trackedEvent[evt.ID.String()] = orderMessage{evt.Sender.String(), msgAction.item, msgAction.provider}
	}
	if err != nil {
		log.Ctx(bot.ctx).Err(err)
		bot.reply(evt.RoomID, evt.ID, err.Error(), false)
		return
	}
	switch len(message) {
	case 0:
		bot.react(evt.RoomID, evt.ID, "✅")
		log.Ctx(bot.ctx).Debug().Msgf("handled actions for message '%s'", evt.Content.AsMessage().Body)
	case 1:
		asHtml := false
		if msgAction.verb == Finalize {
			asHtml = true
		}
		bot.reply(evt.RoomID, evt.ID, message[0], asHtml)
	case 2:
		id := bot.reply(evt.RoomID, evt.ID, message[0], true)
		bot.replyInThread(evt.RoomID, id, message[1])
	default:
		panic("more than 2 messages")
	}
}

func (bot *MatrixBot) handleReactionEvent(evt *event.Event) {
	// get message to which the reaction relates
	msgEvt, err := bot.client.GetEvent(evt.RoomID, evt.Content.AsReaction().RelatesTo.EventID)
	if err != nil {
		log.Ctx(bot.ctx).Error().Err(err).Msg("could not get the message the reaction relates to")
		return
	}
	// parsing has to be triggered manually when `GetEvent` is used
	err = msgEvt.Content.ParseRaw(event.EventMessage)
	if err != nil {
		log.Ctx(bot.ctx).Err(err)
		return
	}
	// get the initial action to use it as base
	action, err := bot.parser.convertToAction(msgEvt.Content.AsMessage().Body)
	switch evt.Content.AsReaction().RelatesTo.Key {
	// run the same action again for the user who reacted
	case "📈":
		var message []string
		message, err = bot.runner.runAction(evt.Sender.String(), action)
		if err != nil {
			log.Ctx(bot.ctx).Err(err)
			return
		}
		bot.reply(evt.RoomID, msgEvt.ID, fmt.Sprintf("%s: ✅ %s", evt.Sender.String(), message), false)
		return
	// remove the item, only possible on own actions
	case "📉":
		if msgEvt.Sender != evt.Sender {
			log.Ctx(bot.ctx).Warn().Msgf("user '%s' tried to troll user '%s'", evt.Sender.String(), msgEvt.Sender.String())
			bot.message(evt.RoomID, fmt.Sprintf("hilarious idea @%s:, but no", evt.Sender.String()))
			return
		}
		action.verb = Remove
		var message []string
		message, err = bot.runner.runAction(evt.Sender.String(), action)
		if err != nil {
			log.Ctx(bot.ctx).Err(err)
			return
		}
		bot.reply(evt.RoomID, msgEvt.ID, fmt.Sprintf("%s: ✅ %s", evt.Sender.String(), message), false)
		return
	}
}

func (bot *MatrixBot) message(room id.RoomID, content string) {
	_, err := bot.client.SendNotice(room, content)
	if err != nil {
		log.Ctx(bot.ctx).Error().Err(err).Msgf("could not send message '%s'", content)
	}
}

func (bot *MatrixBot) react(room id.RoomID, evt id.EventID, content string) {
	_, err := bot.client.SendReaction(room, evt, content)
	if err != nil {
		log.Ctx(bot.ctx).Error().Err(err).Msg("could not react to event")
	}
}

func (bot *MatrixBot) replyInThread(room id.RoomID, evt id.EventID, content string) id.EventID {
	ev, err := bot.client.SendMessageEvent(room, event.EventMessage, map[string]interface{}{
		// "m.mentions": map[string]interface{}{},
		"m.relates_to": map[string]interface{}{
			"rel_type": "m.thread",
			"event_id": evt,
		},
		// notice is a message from a bot, it avoids feedback loops
		"msgtype": "m.notice",
		"body":    content,
	})
	if err != nil {
		log.Ctx(bot.ctx).Error().Err(err).Msgf("could not respond '%s' to event", content)
	}
	return ev.EventID
}

func (bot *MatrixBot) reply(room id.RoomID, evt id.EventID, content string, asHtml bool) id.EventID {
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
		contentJSON["formatted_body"] = "<code>" + strings.TrimSuffix(content, "\n") + "</code>"
	}
	ev, err := bot.client.SendMessageEvent(room, event.EventMessage, contentJSON)
	if err != nil {
		log.Ctx(bot.ctx).Error().Err(err).Msgf("could not respond '%s' to event", content)
	}
	return ev.EventID
}
