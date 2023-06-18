use super::{menu::MenuItem, user::User};

pub enum State {
    Idle,
    TakeOrders,
    Ordered,
}

impl State {
    pub fn do_transition(&self, t: Transition) -> Self {
        match self {
            State::Idle => match t {
                Transition::Help => State::Idle,
                _ => State::Idle
            },
            _ => State::Idle
        }
    }
    
}

pub enum Transition {
    StartOrder,
    AddItem { user: User, item: MenuItem },
    Finalize,
    Cancel,
    Arrived,
    Help,
}

pub trait StateStore {}
