package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
)

type MessageParser struct {
	ctx            context.Context
	trigger        string
	orderProviders []string
}

func NewMessageParser(ctx context.Context, trigger string) *MessageParser {
	req, err := http.NewRequest(http.MethodOptions, ctx.Value(OmegaStarURLKey).(string), nil)
	if err != nil {
		log.Ctx(ctx).Fatal().Err(err)
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Ctx(ctx).Fatal().Err(err)
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		log.Ctx(ctx).Fatal().Err(err)
	}
	var providers []string
	err = json.Unmarshal(b, &providers)
	if err != nil {
		log.Ctx(ctx).Fatal().Err(err)
	}
	return &MessageParser{
		trigger:        trigger,
		orderProviders: providers,
		ctx:            ctx,
	}
}

type ParsedAction struct {
	provider string
	verb     string
	item     string
}

func (parser *MessageParser) convertToAction(msg string) (*ParsedAction, error) {
	msg, found := strings.CutPrefix(msg, parser.trigger)
	if !found {
		return nil, fmt.Errorf("message '%s' does not contain trigger prefix", msg)
	}
	segments := strings.Split(msg, " ")
	if len(segments) < 2 || len(segments) > 3 {
		return nil, fmt.Errorf("message '%s' does not have the expected 2 to 3 segments", msg)
	}
	if !parser.checkOrderProvider(segments[0]) {
		return nil, fmt.Errorf("unknown order provider '%s'", segments[0])
	}
	action := ParsedAction{provider: segments[0]}
	if !parser.checkActionVerb(segments[1]) {
		return nil, fmt.Errorf("invalid action verb '%s'", segments[1])
	}
	action.verb = segments[1]
	if len(segments) == 3 {
		action.item = segments[2]
	}
	return &action, nil
}

func (parser *MessageParser) checkOrderProvider(val string) bool {
	for _, provider := range parser.orderProviders {
		if val == provider {
			return true
		}
	}
	return false
}

func (parser *MessageParser) checkActionVerb(val string) bool {
	for _, verb := range getSupportedActionVerbs() {
		if val == verb {
			return true
		}
	}
	return false
}
