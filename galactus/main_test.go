package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestNewOrderApi(t *testing.T) {
	s := NewServer()
	req, err := http.NewRequest("POST", "/new", nil)
	if err != nil {
		t.Errorf("error when creating request: %s", err.Error())
	}
	rr := httptest.NewRecorder()
	s.newOrder(rr, req)
	if s.nextId != 2 {
		t.Errorf("expected next id to be 2 but is %d", s.nextId)
	}
	if len(s.activeOrders) != 1 {
		t.Errorf("expected one active order but has %d", len(s.activeOrders))
	}
	b, err := io.ReadAll(rr.Result().Body)
	if err != nil {
		t.Errorf("could not read response body: %s", err.Error())
	}
	var newIdMap map[string]int
	err = json.Unmarshal(b, &newIdMap)
	if err != nil {
		t.Errorf("could not unmarshal response: %s", err.Error())
	}
	if newIdMap["id"] != 1 {
		t.Errorf("expected id to be 1 but is %d", newIdMap["id"])
	}
}

func TestStatusApi(t *testing.T) {
	s := NewServer()
	s.activeOrders = make(map[int]*OrderHandler)
	s.activeOrders[1] = NewOrderHandler()
	s.activeOrders[1].orders["user1"] = []string{"M1"}
	req, err := http.NewRequest("GET", "/1/status", nil)
	if err != nil {
		t.Errorf("error when creating request: %s", err.Error())
	}
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/{orderNo}/status", s.orderStatus)
	router.ServeHTTP(rr, req)
	if rr.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected response to be ok but is %s", rr.Result().Status)
	}
	b, err := io.ReadAll(rr.Result().Body)
	defer rr.Result().Body.Close()
	if err != nil {
		t.Fatalf("could not read result body: %s", err.Error())
	}
	var orders map[string][]string
	err = json.Unmarshal(b, &orders)
	if err != nil {
		t.Log(string(b))
		t.Fatalf("error when unmarshaling result: %s", err.Error())
	}
	if _, ok := orders["user1"]; !ok {
		t.Fatal("result body does not contain etries for user1")
	}
	if val := orders["user1"][0]; val != "M1" {
		t.Fatalf("expected order of user1 to be M1 but is %s", val)
	}
}

func TestAddOrder(t *testing.T) {
	s := NewServer()
	s.activeOrders = make(map[int]*OrderHandler)
	s.activeOrders[1] = NewOrderHandler()
	bodyStruct := UpdateOrderBody{
		User: "user1",
		Item: "M1",
	}
	body, err := json.Marshal(&bodyStruct)
	if err != nil {
		t.Errorf("error when creating request: %s", err.Error())
	}
	t.Log(string(body))
	req, err := http.NewRequest("POST", "/1/add", bytes.NewReader(body))
	if err != nil {
		t.Errorf("error when creating request: %s", err.Error())
	}
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/{orderNo}/{action}", s.updateOrder)
	router.ServeHTTP(rr, req)
	if rr.Result().StatusCode != http.StatusCreated {
		defer rr.Result().Body.Close()
		if b, err := io.ReadAll(rr.Result().Body); err != nil {
			t.Log(err.Error())
		} else {
			t.Log(string(b))
		}
		t.Fatalf("expected response to be ok but is %s", rr.Result().Status)
	}
	if _, ok := s.activeOrders[1]; !ok {
		t.Fatal("expected server to have an active order with id 1")
	}
	if _, ok := s.activeOrders[1].orders["user1"]; !ok {
		t.Fatal("expected order with id 1 to have entries for user1")
	}
	if len(s.activeOrders[1].orders["user1"]) != 1 {
		t.Errorf("expected length of orders of user1 to 1 but is: %d", len(s.activeOrders[0].orders["user1"]))
	}
	if s.activeOrders[1].orders["user1"][0] != "M1" {
		t.Errorf("expected order of user1 to be M1 but is: %s", s.activeOrders[0].orders["user1"][0])
	}
}
