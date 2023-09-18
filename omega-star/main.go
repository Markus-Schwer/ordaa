package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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
	CheckItems([]string) []string
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

func (server *MenuServer) start(ctx context.Context) {
	router := mux.NewRouter()
	router.HandleFunc("/", server.getProviderNames).Methods(http.MethodOptions)
	router.HandleFunc("/{provider}", server.getMenu).Methods(http.MethodGet)
	router.HandleFunc("/{provider}/check", server.checkItems).Methods(http.MethodPost)
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
	go server.updateCache(ctx)
	// Run our server in a goroutine so that it doesn't block.
	go func() {
		log.Printf("starting server on %s:%s\n", address, port)
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
	<-ctx.Done()
	// create context to wait for deadline before shutdown
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	log.Println("shutting down")
	srv.Shutdown(ctx)
}

func (server *MenuServer) updateCache(ctx context.Context) {
	ticker := time.NewTicker(12 * time.Hour)
	defer ticker.Stop()
	for _, p := range server.providers {
		log.Printf("updating cache for %s", p.GetName())
		p.UpdateCache()
	}
	for {
		select {
		case <-ctx.Done():
			return
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

func (server *MenuServer) checkItems(w http.ResponseWriter, r *http.Request) {
	name := strings.ToLower(mux.Vars(r)["provider"])
	log.Printf("CHECK %s", name)
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

func main() {
	server := NewMenuServer()
	log.Println("omega star is starting")
	ctx, cancel := context.WithCancel(context.Background())
	go server.start(ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	log.Println("shutting down")
	cancel()
	os.Exit(0)
}
