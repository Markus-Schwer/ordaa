package orders

import "testing"

func TestHappyPath(t *testing.T) {
	oh := newOrderHandler()
	if err := oh.addItem("user1", "M2"); err != nil {
		t.Errorf("could not add item: %s", err.Error())
	}
	if err := oh.addItem("user1", "M1"); err != nil {
		t.Errorf("could not add item: %s", err.Error())
	}
	if err := oh.addItem("user2", "M1"); err != nil {
		t.Errorf("could not add item: %s", err.Error())
	}
	orders, err := oh.finalize()
	if err != nil {
		t.Errorf("could not finalize orders: %s", err.Error())
	}
	for user := range orders {
		t.Log(user)
	}
	if user1orders, ok := orders["user1"]; ok {
		if len(user1orders) != 2 {
			t.Errorf("expected user1 to have 2 orders but has %d", len(user1orders))
		}
	} else {
		t.Error("user1 does not have any orders")
	}
	if user2orders, ok := orders["user2"]; ok {
		if len(user2orders) != 1 {
			t.Errorf("expected user2 to have 1 orders but has %d", len(user2orders))
		}
	} else {
		t.Error("user2 does not have any orders")
	}
	if err := oh.arrived(); err != nil {
		t.Errorf("could not set order to arrived: %s", err.Error())
	}
}

func TestDeleteItem(t *testing.T) {
	oh := newOrderHandler()
	if err := oh.addItem("user1", "M1"); err != nil {
		t.Errorf("could not add item: %s", err.Error())
	}
	allOrders := oh.getOrders()
	if user1orders, ok := allOrders["user1"]; ok {
		if len(user1orders) != 1 {
			t.Errorf("expected user1 to have 1 orders but has %d", len(user1orders))
		}
	}
	if err := oh.removeItem("user1", "M1"); err != nil {
		t.Errorf("could not remove item: %s", err.Error())
	}
	allOrders = oh.getOrders()
	if user1orders, ok := allOrders["user1"]; ok {
		if len(user1orders) != 0 {
			t.Errorf("expected user1 to have 1 orders but has %d", len(user1orders))
		}
	}
}

func TestFailOnOrderInDelivering(t *testing.T) {
	oh := newOrderHandler()
	if err := oh.addItem("user1", "M1"); err != nil {
		t.Errorf("could not add item: %s", err.Error())
	}
	if _, err := oh.finalize(); err != nil {
		t.Errorf("could not finalize: %s", err.Error())
	}
	if err := oh.addItem("user2", "M1"); err == nil {
		t.Error("expected order handler to throw error when ordering after finalizing")
	}
}
