use super::{menu::MenuItem, user::User};

#[derive(Clone)]
pub enum State {
    Idle,
    TakeOrders,
    Ordered,
}

impl State {
    pub fn do_transition(&self, t: Action) -> Option<Self> {
    }
}

#[derive(PartialEq)]
pub enum Action {
    StartOrder,
    AddItem { user: User, item: MenuItem },
    Finalize,
    Cancel,
    Arrived,
    Help,
}

pub trait StateStore {}
