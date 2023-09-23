package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/galactus/orders"
)

type RestInterface struct {
	mo            *orders.MultiOrders
	serverContext context.Context
}

func NewRestInterface(mo *orders.MultiOrders) RestInterface {
	return RestInterface{
		mo: mo,
	}
}

func (server *RestInterface) start(ctx context.Context) {
	server.serverContext = ctx
	router := mux.NewRouter()
	router.HandleFunc("/{provider}/new", server.newOrder).Methods(http.MethodPost)
	router.HandleFunc("/{orderNo}/{action}", server.updateOrder).Methods(http.MethodPost)
	router.HandleFunc("/{orderNo}/status", server.orderStatus).Methods(http.MethodGet)
	router.HandleFunc("/status", server.ordersStatus).Methods(http.MethodGet)
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
	<-ctx.Done()
	// create context to wait for deadline before shutdown
	shutdownContext, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	srv.Shutdown(shutdownContext)
}

func (server *RestInterface) newOrder(w http.ResponseWriter, r *http.Request) {
	orderNo := server.mo.CreateNewOrder(mux.Vars(r)["provider"])
	if _, err := fmt.Fprintf(w, "{\"orderNo\": %d}", orderNo); err != nil {
		http.Error(w, fmt.Sprintf("could not write order id to response: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
}

func (server *RestInterface) orderStatus(w http.ResponseWriter, r *http.Request) {
	orderNo, err := strconv.Atoi(mux.Vars(r)["orderNo"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	order, provider, err := server.mo.GetOrder(orderNo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if order == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	b, err := json.Marshal(map[string]interface{}{
		"orders":   order,
		"provider": provider,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("could not serialize orders: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	w.Write(b)
	w.Header().Add("Content-Type", "application/json")
}

func (server *RestInterface) ordersStatus(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(server.mo.GetOrders())
	if err != nil {
		http.Error(w, fmt.Sprintf("could not serialize orders meta: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	w.Write(b)
	w.Header().Add("Content-Type", "application/json")
}

type UpdateOrderBody struct {
	User string `json:"user"`
	Item string `json:"item"`
}

func (server *RestInterface) updateOrder(w http.ResponseWriter, r *http.Request) {
	orderNo, err := strconv.Atoi(mux.Vars(r)["orderNo"])
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid orderNo URL parameter: %s", err.Error()), http.StatusBadRequest)
		return
	}
	action := mux.Vars(r)["action"]
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
	order, err := server.mo.HandleOrderAction(server.serverContext, orders.OrderAction{
		Action:  action,
		User:    data.User,
		Item:    data.Item,
		OrderNo: orderNo,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("could not handle action: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if order == nil {
		w.WriteHeader(http.StatusOK)
		return
	}
	b, err = json.Marshal(order)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not marshal response: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	w.Write(b)
}
