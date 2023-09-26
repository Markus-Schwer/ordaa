package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

func NewGalactusFacade(ctx context.Context, timeout time.Duration) ServiceFacade {
	return &Microfacade{
		ctx:     ctx,
		timeout: timeout,
	}
}

type Microfacade struct {
	ctx     context.Context
	timeout time.Duration
}

func (facade *Microfacade) CheckOrderItem(provider string, item string) error {
	b, err := facade.doHttp(
		fmt.Sprintf("%s/%s/check", facade.ctx.Value(OmegaStarURLKey).(string), provider),
		http.MethodPost,
		[]string{item},
	)
	if err != nil {
		return fmt.Errorf("omga star failed to vaidate: %s", err.Error())
	}
	var invalidItems []string
	err = json.Unmarshal(b, &invalidItems)
	if err != nil {
		return err
	}
	if len(invalidItems) == 0 {
		return nil
	}
	return fmt.Errorf("there are invalid items: %v", invalidItems)
}

func (facade *Microfacade) NewOrder(provider string) (int, error) {
	b, err := facade.doHttp(
		fmt.Sprintf("%s/%s/new", facade.ctx.Value(GalactusURLKey), provider),
		http.MethodPost,
		nil,
	)
	if err != nil {
		return 0, fmt.Errorf("galactus failed to vaidate: %s", err.Error())
	}
	var orderNoMap map[string]int
	err = json.Unmarshal(b, &orderNoMap)
	if err != nil {
		return 0, err
	}
	orderNo, ok := orderNoMap["orderNo"]
	if !ok {
		return 0, fmt.Errorf("omega star returned wird json: %v", orderNoMap)
	}
	return orderNo, nil
}

func (facade *Microfacade) AddOrderItem(orderNo int, user string, item string) error {
	_, err := facade.doHttp(
		fmt.Sprintf("%s/%d/add", facade.ctx.Value(GalactusURLKey), orderNo),
		http.MethodPost,
		map[string]string{
			"user": user,
			"item": item,
		},
	)
	return err
}

func (facade *Microfacade) RemoveOrderItem(orderNo int, user string, item string) error {
	_, err := facade.doHttp(
		fmt.Sprintf("%s/%d/remove", facade.ctx.Value(GalactusURLKey), orderNo),
		http.MethodPost,
		map[string]string{
			"user": user,
			"item": item,
		},
	)
	return err
}

func (facade *Microfacade) FinalizeOrder(orderNo int) (orders Orders, err error) {
	b, err := facade.doHttp(
		fmt.Sprintf("%s/%d/finalize", facade.ctx.Value(GalactusURLKey), orderNo),
		http.MethodPost,
		nil,
	)
	if err != nil {
		return
	}
	err = json.Unmarshal(b, &orders)
	return
}

func (facade *Microfacade) OrderArrived(orderNo int) error {
	_, err := facade.doHttp(
		fmt.Sprintf("%s/%d/arrived", facade.ctx.Value(GalactusURLKey), orderNo),
		http.MethodPost,
		nil,
	)
	return err
}

func (facade *Microfacade) CancelOrder(orderNo int) error {
	_, err := facade.doHttp(
		fmt.Sprintf("%s/%d/cancel", facade.ctx.Value(GalactusURLKey), orderNo),
		http.MethodPost,
		nil,
	)
	return err
}

func (facade *Microfacade) GetOrders() (orderMeta []OrderMetadata, err error) {
	b, err := facade.doHttp(
		fmt.Sprintf("%s/status", facade.ctx.Value(GalactusURLKey)),
		http.MethodGet,
		nil,
	)
	if err != nil {
		return
	}
	log.Ctx(facade.ctx).Debug().Bytes("bytes", b).Msg("order response from galactus arrived")
	err = json.Unmarshal(b, &orderMeta)
	if err != nil {
		return
	}
	return
}

// can handle nil body
func (facade *Microfacade) doHttp(url string, method string, jsonBody interface{}) ([]byte, error) {
	var bodyReader io.Reader
	if jsonBody != nil {
		b, err := json.Marshal(jsonBody)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(b)
	}
	withTo, cancel := context.WithTimeout(facade.ctx, facade.timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(
		withTo,
		method,
		url,
		bodyReader,
	)
	if jsonBody != nil {
		req.Header.Add("Content-Type", "application/json")
	}
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	b, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		if err != nil {
			return nil, err
		}
		return b, fmt.Errorf("%s", string(b))
	}
	return b, nil
}
