package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

type MenuItem struct {
	Id    string
	Name  string
	Price float32
}

type Menu struct {
	Items []MenuItem
}

type MenuProvider interface {
	GetName() string
	UpdateCache() error
	GetMenu() *Menu
}

type MenuServer struct {
	providers []MenuProvider
}

func NewMenuServer() MenuServer {
	return MenuServer{
		providers: []MenuProvider{
			InitSangam(),
		},
	}
}

func (server *MenuServer) start() {
	router := mux.NewRouter()
	router.HandleFunc("/", server.getProviderNames).Methods(http.MethodOptions)
	router.HandleFunc("/{provider}", server.getMenu).Methods(http.MethodGet)
	address := os.Getenv("OMEGA_STAR_ADDRESS")
	if address == "" {
		address = "127.0.0.1"
	}
	port := os.Getenv("OMEGA_STAR_PORT")
	if port == "" {
		port = "8080"
	}
	srv := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf("%s:%s", address, port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	ticker := time.NewTicker(12 * time.Hour)
	go server.updateCache(ticker)
	defer ticker.Stop()
	// Run our server in a goroutine so that it doesn't block.
	go func() {
		log.Printf("starting server on %s:%s\n", address, port)
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
	// channel to notify when shutdown should happen
	c := make(chan os.Signal, 1)
	// graceful shutdown on SIGINT and SIGTERM
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	// block until shutdown is desired
	<-c
	// create context to wait for deadline before shutdown
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}

func (server *MenuServer) updateCache(ticker *time.Ticker) {
	for _, p := range server.providers {
		log.Printf("updating cache for %s", p.GetName())
		p.UpdateCache()
	}
	for {
		select {
		case <-ticker.C:
			log.Println("updating cache")
			for _, p := range server.providers {
				if err := p.UpdateCache(); err != nil {
					log.Println(err)
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
	log.Printf("GET %s", name)
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

func main() {
	server := NewMenuServer()
	server.start()
}
