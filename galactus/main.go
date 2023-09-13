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
	"strconv"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	s := NewServer()
	s.start()
}

func NewServer() Server {
	return Server{
		activeOrders: make(map[int]*OrderHandler),
		nextId:       1,
	}
}

type Server struct {
	activeOrders map[int]*OrderHandler
	nextId       int
}

func (server *Server) start() {
	router := mux.NewRouter()
	router.HandleFunc("/new", server.newOrder).Methods(http.MethodPost)
	router.HandleFunc("/{orderNo}/{action}", server.updateOrder).Methods(http.MethodPost)
	router.HandleFunc("/{orderNo}/status", server.orderStatus).Methods(http.MethodGet)
	address := os.Getenv("GALACTUS_ADDRESS")
	if address == "" {
		address = "127.0.0.1"
	}
	port := os.Getenv("GALACTUS_PORT")
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

func (server *Server) newOrder(w http.ResponseWriter, r *http.Request) {
	log.Printf("create new order with id %d", server.nextId)
	server.activeOrders[server.nextId] = NewOrderHandler()
	retVal := map[string]int{"id": server.nextId}
	server.nextId += 1
	b, err := json.Marshal(retVal)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not marshal order id: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	w.Write(b)
	w.Header().Add("Content-Type", "application/json")
}

func (server *Server) orderStatus(w http.ResponseWriter, r *http.Request) {
	orderNo, err := strconv.Atoi(mux.Vars(r)["orderNo"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("request status of order id %d", orderNo)
	om, ok := server.activeOrders[orderNo]
	if !ok {
		http.Error(w, fmt.Sprintf("no orders for order with id %d", orderNo), http.StatusNotFound)
		return
	}
	b, err := json.Marshal(om.getOrders())
	if err != nil {
		http.Error(w, fmt.Sprintf("could not serialize orders: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	w.Write(b)
	w.Header().Add("Content-Type", "application/json")
}

type UpdateOrderBody struct {
	User string `json:"user"`
	Item string `json:"item"`
}

func (server *Server) updateOrder(w http.ResponseWriter, r *http.Request) {
	orderNo, err := strconv.Atoi(mux.Vars(r)["orderNo"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("update order id %d", orderNo)
	om, ok := server.activeOrders[orderNo]
	if !ok {
		http.Error(w, fmt.Sprintf("order number %d does not exist", orderNo), http.StatusNotFound)
		return
	}
	action := mux.Vars(r)["action"]
	switch action {
	case "add":
		var data UpdateOrderBody
		b, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		if err := json.Unmarshal(b, &data); err != nil {
			http.Error(w, fmt.Sprintf("could not unmarshal json body: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		if err := om.addItem(data.User, data.Item); err != nil {
			http.Error(w, fmt.Sprintf("could not add item to order: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		return
	case "remove":
		var data *UpdateOrderBody
		b, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		json.Unmarshal(b, data)
		if err := om.removeItem(data.User, data.Item); err != nil {
			http.Error(w, fmt.Sprintf("could not remove item from order: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		return
	case "finalize":
		b, err := json.Marshal(om.getOrders())
		if err != nil {
			http.Error(w, fmt.Sprintf("could not finalize order: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		w.Write(b)
		w.Header().Add("Content-Type", "application/json")
		return
	case "arrived":
		if err := om.arrived(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	case "cancel":
		if err := om.cancel(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	case "":
		http.Error(w, "invalid empty action", http.StatusBadRequest)
		return
	default:
		http.Error(w, fmt.Sprintf("unknown action: %s", action), http.StatusBadRequest)
		return
	}
}
