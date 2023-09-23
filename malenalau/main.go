package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	HomeServerURLKey = "homeServerURL"
	PasswordKey      = "password"
	UserKey          = "user"
	OmegaStarURLKey  = "omegaStarURL"
	GalactusURLKey   = "GalactusURL"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	ctx = context.WithValue(ctx, HomeServerURLKey, "matrix.org")
	ctx = context.WithValue(ctx, PasswordKey, "2sm^ziai#95nRHoMVzKz")
	ctx = context.WithValue(ctx, UserKey, "order-bot-aa")
	ctx = context.WithValue(ctx, OmegaStarURLKey, "http://localhost:8081")
	ctx = context.WithValue(ctx, GalactusURLKey, "http://localhost:8080")

	facade := NewGalactusFacade(ctx, time.Second*10)
	parser := NewMessageParser(ctx, ".")
	runner := NewActionRunner(facade)
	bot := NewMatrixBot(ctx, parser, runner)

	bot.LoginAndJoin([]string{"!hvJGXMkjcyzxtSNNsx:matrix.org"})
	go bot.Listen()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	cancel()
}
