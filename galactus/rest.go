package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/galactus/orders"
)

type ResponseSignalWrapper struct {
	w    http.ResponseWriter
	done chan<- bool
}

type RestInterface struct {
	in  chan<- orders.OrderAction
	out chan orders.OrderActionResponse
	ctx context.Context
	// map of uuid to response writers
	responseMap sync.Map
}

func NewRestInterface(ctx context.Context, in chan<- orders.OrderAction, out chan orders.OrderActionResponse) RestInterface {
	return RestInterface{
		ctx:         ctx,
		in:          in,
		out:         out,
		responseMap: sync.Map{},
	}
}

func (server *RestInterface) start() {
	router := mux.NewRouter()
	router.HandleFunc("/{provider}/new", server.newOrder).Methods(http.MethodPost)
	router.HandleFunc("/{orderNo}/{action}", server.updateOrder).Methods(http.MethodPost)
	router.HandleFunc("/{orderNo}/status", server.orderStatus).Methods(http.MethodGet)
	router.HandleFunc("/status", server.ordersStatus).Methods(http.MethodGet)
	router.Use(server.loggingMiddleware)
	srv := &http.Server{
		Handler:      router,
		Addr:         server.ctx.Value(AddressKey).(string),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Ctx(server.ctx).Error().Err(err).Msg("listen and serve crashed")
		}
	}()
	for {
		select {
		case res := <-server.out:
			wrapper, ok := server.responseMap.LoadAndDelete(res.Uuid())
			if !ok {
				log.Ctx(server.ctx).Debug().Msgf("got response for %s but found no writer", res.Uuid())
				break
			}
			v, ok := wrapper.(ResponseSignalWrapper)
			if !ok {
				log.Ctx(server.ctx).Fatal().Msg("stored value in map was no response signal wrapper")
				break
			}
			go server.writeResponse(res, v)
		case <-server.ctx.Done():
			log.Ctx(server.ctx).Debug().Msg("starting 5 second grace period on shutdown")
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()
			srv.Shutdown(ctx)
			log.Ctx(server.ctx).Debug().Msg("HTTP server shutdown complete")
			return
		}
	}
}

func (server *RestInterface) writeResponse(genericResponse orders.OrderActionResponse, wrapper ResponseSignalWrapper) {
	switch res := genericResponse.(type) {
	case *orders.OkWithOrderNo:
		if _, err := fmt.Fprintf(wrapper.w, "{\"orderNo\": %d}", res.OrderNo); err != nil {
			log.Ctx(server.ctx).Error().Err(err).Msg("could not write response")
			http.Error(wrapper.w, "I can't write", http.StatusInternalServerError)
		}
		wrapper.w.Header().Add("Content-Type", "application/json")
	case *orders.Ok:
		wrapper.w.WriteHeader(http.StatusOK)
	case *orders.GenericError:
		http.Error(wrapper.w, res.Error().Error(), http.StatusInternalServerError)
	case *orders.NoActiveOrder:
		http.Error(wrapper.w, "did you miss the order? too bad", http.StatusBadRequest)
	case *orders.OkWithOrder:
		b, err := json.Marshal(res.Order)
		if err != nil {
			http.Error(wrapper.w, "I JSON fumbled", http.StatusInternalServerError)
		}
		if _, err := wrapper.w.Write(b); err != nil {
			log.Ctx(server.ctx).Error().Err(err).Msg("could not write response")
			http.Error(wrapper.w, "I can't write", http.StatusInternalServerError)
		}
		wrapper.w.Header().Add("Content-Type", "application/json")
	case *orders.OkWithActiveOrders:
		b, err := json.Marshal(res.ActiveOrders)
		if err != nil {
			http.Error(wrapper.w, "I JSON fumbled", http.StatusInternalServerError)
		}
		if _, err := wrapper.w.Write(b); err != nil {
			log.Ctx(server.ctx).Error().Err(err).Msg("could not write response")
			http.Error(wrapper.w, "I can't write", http.StatusInternalServerError)
		}
		wrapper.w.Header().Add("Content-Type", "application/json")
	default:
		log.Ctx(server.ctx).Error().Msg("unhandled response type")
		http.Error(wrapper.w, "I've never heard such a response", http.StatusInternalServerError)
	}
	wrapper.done <- false
}

func (server *RestInterface) newOrder(w http.ResponseWriter, r *http.Request) {
	actionUuid := uuid.New()
	server.in <- orders.OrderAction{Action: "new", Provider: mux.Vars(r)["provider"], Uuid: actionUuid}
	done := make(chan bool, 1)
	server.responseMap.Store(actionUuid, ResponseSignalWrapper{w, done})
	<-done
}

func (server *RestInterface) orderStatus(w http.ResponseWriter, r *http.Request) {
	orderNo, err := strconv.Atoi(mux.Vars(r)["orderNo"])
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg("'orderNo' URL segment is not an int")
		http.Error(w, "Use the numba of the order", http.StatusBadRequest)
		return
	}
	actionUuid := uuid.New()
	server.in <- orders.OrderAction{
		Uuid:    actionUuid,
		Action:  "status",
		OrderNo: orderNo,
	}
	done := make(chan bool, 1)
	server.responseMap.Store(actionUuid, ResponseSignalWrapper{w, done})
	<-done
}

func (server *RestInterface) ordersStatus(w http.ResponseWriter, r *http.Request) {
	actionUuid := uuid.New()
	server.in <- orders.OrderAction{
		Uuid:   actionUuid,
		Action: "active",
	}
	done := make(chan bool, 1)
	server.responseMap.Store(actionUuid, ResponseSignalWrapper{w, done})
	<-done
}

type UpdateOrderBody struct {
	User string `json:"user"`
	Item string `json:"item"`
}

func (server *RestInterface) updateOrder(w http.ResponseWriter, r *http.Request) {
	orderNo, err := strconv.Atoi(mux.Vars(r)["orderNo"])
	if err != nil {
		log.Ctx(server.ctx).Error().Err(err).Msg("'orderNo' URL segment is not an int")
		http.Error(w, "Use the numba of the order", http.StatusBadRequest)
		return
	}
	actionUuid := uuid.New()
	action := orders.OrderAction{
		Action:  mux.Vars(r)["action"],
		Uuid:    actionUuid,
		OrderNo: orderNo,
	}
	if action.Action == "add" || action.Action == "remove" {
		var data UpdateOrderBody
		b, err := io.ReadAll(r.Body)
		if err != nil {
			log.Ctx(server.ctx).Error().Err(err).Msg("could not read request body")
			http.Error(w, "I can't read", http.StatusBadRequest)
			return
		}
		if err := json.Unmarshal(b, &data); err != nil {
			log.Ctx(server.ctx).Error().Err(err).Msg("could not parse JSON request body")
			http.Error(w, "This is supposed to be JSON?", http.StatusBadRequest)
			return
		}
		action.Item = data.Item
		action.User = data.User
	}
	server.in <- action
	done := make(chan bool, 1)
	server.responseMap.Store(actionUuid, ResponseSignalWrapper{w, done})
	<-done
}

func (server *RestInterface) loggingMiddleware(next http.Handler) http.Handler {
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
