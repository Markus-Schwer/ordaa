use crate::{users::User, menu::MenuItem};

use super::state::State;

pub struct Order {
    // scoped to orders
    id: i64,
    // time when the order will automatically finalized
    order_deadline: Option<chrono::DateTime<chrono::Local>>,
    // estimated eta
    eta: Option<chrono::DateTime<chrono::Local>>,
    items: Vec<OrderItem>,
    // person who is responisble for the order
    initiator: User,
    // person who paid for the order
    sugar_person: Option<User>,
    state: State
}

pub struct OrderItem {
    id: i64,
    user: User,
    item: MenuItem,
    paid: bool
}

impl Order {
    // pub fn init(dto: CreateOrderDTO, id: i64) -> Self {
    //     Order { 
    //         id,
    //         order_deadline: dto.order_deadline,
    //         eta: None,
    //         items: Vec::new(),
    //         initiator: dto.initiator,
    //         sugar_person: None,
    //         state: State::PlacingOrder
    //     }
    // }

    pub fn add_item(&mut self, it: OrderItem) -> Result<(), String> {
        if !matches!(self.state, State::CollectingOrderItems) {
            return Err(format!("cannot add items to order in state {}", self.state));
        }
        self.items.push(it);
        Ok(())
    }

    pub fn remove_item(&mut self, id: i64) -> Result<(), String> {
        if !matches!(self.state, State::CollectingOrderItems) {
            return Err(format!("cannot remove items from order in state {}", self.state));
        }
        let size_before = self.items.len();
        self.items.retain(|it| it.id != id);
        if self.items.len() == size_before {
            return Err(format!("item with id {} was not in the order", id));
        }
        Ok(())
    }

    pub fn finalize_order(&mut self) -> Result<(), String> {
        match self.state {
            State::CollectingOrderItems => {},
            State::PlacingOrder => return Err("that isn't necessary, it's already locked in".to_string()),
            State::InDelivery => return Err("the cook's already on it, and guess whose order is missing".to_string()),
            State::Arrived => return Err("too bad, the others are already eating or maybe even done".to_string()),
        }
        self.state = State::PlacingOrder;
        Ok(())
    }

    pub fn in_delivery(&mut self, eta: Option<chrono::DateTime<chrono::Local>>) -> Result<(), String> {
        match self.state {
            State::CollectingOrderItems => return Err("I hope you didn't really order anything, we're still collecting".to_string()),
            State::PlacingOrder => {},
            State::InDelivery => return Err("that isn't necessary, the cook's already on it".to_string()),
            State::Arrived => return Err("too bad, the others are already eating or maybe even done".to_string()),
        }
        self.state = State::InDelivery;
        self.eta = eta;
        Ok(())
    }

    pub fn arrived(&mut self) -> Result<(), String> {
        match self.state {
            State::CollectingOrderItems => return Err("what arrived? we wouldn't even know, we haven't decided yet".to_string()),
            State::PlacingOrder => return Err(format!("did {} conjoure it from nothing or how did it arrive before even ordering?", self.initiator)),
            State::InDelivery => {},
            State::Arrived => return Err("you probably stole someone else's order, because this one already arrived".to_string()),
        }
        self.state = State::Arrived;
        Ok(())
    }
}
