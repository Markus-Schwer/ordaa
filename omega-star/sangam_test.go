package main

import (
	"os"
	"testing"
)

func TestParsing(t *testing.T) {
	s := InitSangam()
	b, err := os.ReadFile("test-resources/sangam.html")
	if err != nil {
		t.Fatal(err.Error())
	}
	s.cachedHtml = string(b)
	err = s.updateMenuFromCache()
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(s.GetMenu().Items) != 136 {
		t.Errorf("expected %d items on the menu but was %d", 136, len(s.GetMenu().Items))
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
