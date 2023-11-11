package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type MenuItem struct {
	Id    string
	Name  string
	Price int
}

type Menu struct {
	Items []MenuItem
}

type MenuProvider interface {
	GetName() string
	UpdateCache() error
	GetMenu() *Menu
	CheckItems([]string) []string
}

type MenuServer struct {
	ctx       context.Context
	providers []MenuProvider
}

func NewMenuServer(ctx context.Context) MenuServer {
	return MenuServer{
		ctx: ctx,
		providers: []MenuProvider{
			InitSangam(ctx),
		},
	}
}

func (server *MenuServer) start() {
	router := mux.NewRouter()
	router.HandleFunc("/", server.getProviderNames).Methods(http.MethodOptions)
	router.HandleFunc("/{provider}", server.getMenu).Methods(http.MethodGet)
	router.HandleFunc("/{provider}/check", server.checkItems).Methods(http.MethodPost)
	router.Use(server.loggingMiddleware)
	srv := &http.Server{
		Handler: router,
		Addr:    server.ctx.Value(AddressKey).(string),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	go server.updateCache()
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Ctx(server.ctx).Error().Err(err).Msg("listen and serve crashed")
		}
	}()
	<-server.ctx.Done()
	log.Ctx(server.ctx).Debug().Msg("starting 5 second grace period on shutdown")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	srv.Shutdown(ctx)
	log.Ctx(server.ctx).Debug().Msg("HTTP server shutdown complete")
}

func (server *MenuServer) updateCache() {
	ticker := time.NewTicker(12 * time.Hour)
	defer ticker.Stop()
	for _, p := range server.providers {
		p.UpdateCache()
	}
	for {
		select {
		case <-server.ctx.Done():
			log.Ctx(server.ctx).Debug().Msg("cache updater received shutdown signal")
			return
		case <-ticker.C:
			log.Ctx(server.ctx).Info().Msg("updating cache")
			for _, p := range server.providers {
				if err := p.UpdateCache(); err != nil {
					log.Ctx(server.ctx).Error().Err(err).Msg("updating cache failed")
				}
			}
		}
	}
}

func (server *MenuServer) getProviderNames(w http.ResponseWriter, r *http.Request) {
	options := make([]string, 0)
	for _, p := range server.providers {
		options = append(options, strings.ToLower(p.GetName()))
	}
	json.NewEncoder(w).Encode(options)
}

func (server *MenuServer) getMenu(w http.ResponseWriter, r *http.Request) {
	name := strings.ToLower(mux.Vars(r)["provider"])
	for _, provider := range server.providers {
		if provider.GetName() == name {
			if menu := provider.GetMenu(); menu != nil {
				json.NewEncoder(w).Encode(provider.GetMenu())
				return
			} else {
				http.Error(w, "", http.StatusServiceUnavailable)
				return
			}
		}
	}
	http.Error(w, "", http.StatusNotFound)
}

func (server *MenuServer) checkItems(w http.ResponseWriter, r *http.Request) {
	name := strings.ToLower(mux.Vars(r)["provider"])
	var data []string
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not read body: %s", err.Error()), http.StatusInternalServerError)
	}
	err = json.Unmarshal(b, &data)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not read body: %s", err.Error()), http.StatusInternalServerError)
	}
	for _, provider := range server.providers {
		if provider.GetName() != name {
			continue
		}
		invalid := provider.CheckItems(data)
		b, err := json.Marshal(invalid)
		if err != nil {
			http.Error(w, fmt.Sprintf("could not write response: %s", err.Error()), http.StatusInternalServerError)
		}
		w.Write(b)
		return
	}
	http.Error(w, "", http.StatusNotFound)
}

func (server *MenuServer) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		respWriterSpy := NewResponseWriterSpy(w)
		next.ServeHTTP(respWriterSpy, r)
		log.Ctx(server.ctx).Info().
			Stringer("route", r.URL).
			Str("user-agent", r.UserAgent()).
			Str("response-code", fmt.Sprintf("%d", respWriterSpy.statusCode)).
			Stringer("duration", time.Since(startTime)).
			Msg("served request")
	})
}
