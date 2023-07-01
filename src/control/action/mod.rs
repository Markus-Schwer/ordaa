use enum_dispatch::enum_dispatch;

use super::state::State;

pub mod add_item;
pub mod arrived;
pub mod cancel;
pub mod finalize;
pub mod help;
pub mod start_order;

use add_item::AddItem;
use arrived::Arrived;
use cancel::Cancel;
use finalize::Finalize;
use help::Help;
use start_order::StartOrder;

#[derive(PartialEq, Eq, Hash)]
#[enum_dispatch]
pub enum ActionEnum {
    StartOrder,
    AddItem,
    Finalize,
    Cancel,
    Arrived,
    Help,
}

#[enum_dispatch(ActionEnum)]
pub trait Action {
    // it feels bad to copy the state every call, but passing it mutably is
    // just as bad.
    // TODO: optimize this if possible
    fn reduce(&self, state: State) -> Result<State, ReducerError>;
}

#[derive(Debug)]
pub enum ReducerError {
    InvalidTransition { message: String },
    InvalidState { message: String },
}
