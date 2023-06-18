use super::{menu::MenuItem, user::User};

#[derive(Clone)]
pub enum State {
    Idle,
    TakeOrders,
    Ordered,
}

impl State {
    pub fn do_transition(&self, t: Transition) -> Option<Self> {
        if t == Transition::Help {
            return Some((*self).clone());
        }
        match self {
            State::Idle => match t {
                Transition::StartOrder => Some(State::TakeOrders),
                _ => None,
            },
            State::TakeOrders => match t {
                Transition::AddItem { user: _, item: _ } => Some(State::TakeOrders),
                Transition::Finalize => Some(State::Ordered),
                _ => None,
            },
            State::Ordered => match t {
                Transition::Arrived => Some(State::Idle),
                _ => None,
            },
        }
    }
}

#[derive(PartialEq)]
pub enum Transition {
    StartOrder,
    AddItem { user: User, item: MenuItem },
    Finalize,
    Cancel,
    Arrived,
    Help,
}

pub trait StateStore {}
