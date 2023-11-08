package main

import (
	"context"
	"os"
	"testing"
)

func TestParsing(t *testing.T) {
	s := InitSangam(context.Background())
	b, err := os.ReadFile("test-resources/sangam.html")
	if err != nil {
		t.Fatal(err.Error())
	}
	s.cachedHtml = string(b)
	err = s.updateMenuFromCache()
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(s.GetMenu().Items) != 134 {
		t.Errorf("expected %d items on the menu but was %d", 134, len(s.GetMenu().Items))
	}
	for _, item := range s.GetMenu().Items {
		if item.Id == "" {
			t.Errorf("item %s has an empty id", item.Name)
		}
		if item.Name == "" {
			t.Errorf("item %s has an empty name", item.Id)
		}
		if item.Price == 0 {
			t.Errorf("item %s has an empty price", item.Name)
		}
	}
}

func TestChecks(t *testing.T) {
	s := InitSangam(context.Background())
	// s.menu = &Menu{Items: make([]MenuItem, 0, len(*s.menu))}
    newMenu := make(map[string]MenuItem)
    newMenu["M1"] = MenuItem{Id: "M1", Name: "Menu 1", Price: 60}
    newMenu["M2"] = MenuItem{Id: "M2", Name: "Menu 2", Price: 420}
	s.menu = &newMenu
	ret := s.CheckItems([]string{"M1", "M3"})
	if len(ret) != 1 {
		t.Fatalf("expected to find 1 invalid item but got %d", len(ret))
	}
	if ret[0] != "M3" {
		t.Fatal("expected first invalid item to be 'M3'")
	}
}
