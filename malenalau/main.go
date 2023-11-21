package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	HomeServerURLKey = "homeServerURL"
	PasswordKey      = "password"
	UserKey          = "user"
	OmegaStarURLKey  = "omegaStarURL"
	GalactusURLKey   = "GalactusURL"
)

func main() {
	var verbose, jsonFormat bool
	flag.BoolVar(&verbose, "v", false, "verbose output: sets the log level to debug")
	flag.BoolVar(&jsonFormat, "j", false, "logging in json format")
	flag.Parse()

	room := os.Getenv("ROOM")
	homeServerUrl := os.Getenv("HOME_SERVER")
	user := os.Getenv("USER")
	passwordFile := os.Getenv("PASSWORD_FILE")
	omegaStarUrl := os.Getenv("OMEGA_STAR_URL")
	galactusUrl := os.Getenv("GALACTUS_URL")

	b, err := os.ReadFile(passwordFile)
	if err != nil {
		log.Fatal().Err(err).Msgf("could not read password file '%s'", passwordFile)
	}

	ctx, cancel := context.WithCancel(context.Background())
	ctx = context.WithValue(ctx, HomeServerURLKey, homeServerUrl)
	ctx = context.WithValue(ctx, PasswordKey, strings.TrimSpace(string(b)))
	ctx = context.WithValue(ctx, UserKey, user)
	ctx = context.WithValue(ctx, OmegaStarURLKey, omegaStarUrl)
	ctx = context.WithValue(ctx, GalactusURLKey, galactusUrl)

	if jsonFormat {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	ctx = log.With().Str("service", "malenalau").Logger().WithContext(ctx)
	if verbose {
		log.Ctx(ctx).Level(zerolog.DebugLevel)
	} else {
		log.Ctx(ctx).Level(zerolog.InfoLevel)
	}
	log.Ctx(ctx).Debug().Msgf("user: '%s'", user)

	facade := NewGalactusFacade(ctx, time.Second*10)
	parser := NewMessageParser(ctx, ".")
	runner := NewActionRunner(ctx, facade)
	bot := NewMatrixBot(ctx, parser, runner)

	bot.LoginAndJoin([]string{room})
	go bot.Listen()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	cancel()
	os.Exit(0)
}
