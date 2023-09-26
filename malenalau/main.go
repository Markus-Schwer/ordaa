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
	var homeServerUrl, user, passwordFile, omegaStarUrl, galactusUrl string
	var rooms RoomsFlag
	var verbose bool
	flag.BoolVar(&verbose, "v", false, "verbose output: sets the log level to debug")
	flag.StringVar(&homeServerUrl, "home-server-url", "matrix.org", "matrix home server url, e.g. matrix.org")
	flag.StringVar(&user, "user", "", "username of the matrix account")
	flag.StringVar(&passwordFile, "password-file", "", "location of the file with the passowrd for the matrix account")
	flag.StringVar(&omegaStarUrl, "omega-star", "http://localhost:8081", "URL where the omega star service can be reached")
	flag.StringVar(&galactusUrl, "galactus", "http://localhost:8080", "URL where the galactus service can be reached")
	flag.Var(&rooms, "room", "repeatable flag with matrix room ids the bot should join")
	flag.Parse()

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

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
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

	bot.LoginAndJoin(rooms)
	go bot.Listen()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	cancel()
	os.Exit(0)
}

type RoomsFlag []string

func (rooms *RoomsFlag) Set(value string) error {
	*rooms = append(*rooms, value)
	return nil
}

func (rooms *RoomsFlag) String() (out string) {
	out = "["
	out += strings.Join(*rooms, ",")
	out += "]"
	return
}
