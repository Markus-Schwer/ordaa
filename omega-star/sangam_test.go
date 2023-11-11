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
		if item.Id == "124" && item.Name != "Karahi-Fisch-Lababer" {
			t.Errorf("item 124 should have name 'Karahi-Fisch-Lababer' but is %s", item.Name)
		}
		if item.Id == "M12" && item.Price != 890 {
			t.Errorf("item M12 should have price '890' but is %d", item.Price)
		}
		if item.Id == "25" && item.Price != 990 {
			t.Errorf("item 25 should have price '990' but is %d", item.Price)
		}
	}
}

func TestChecks(t *testing.T) {
	s := InitSangam(context.Background())
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
