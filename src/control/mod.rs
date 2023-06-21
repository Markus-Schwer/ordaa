use std::collections::HashMap;

use async_trait::async_trait;
use enum_dispatch::enum_dispatch;
use tokio::sync::mpsc::{UnboundedSender, UnboundedReceiver};

use reducer::start_order::StartOrder;
use reducer::add_item::AddItem;
use reducer::finalize::Finalize;
use reducer::cancel::Cancel;
use reducer::arrived::Arrived;
use reducer::help::Help;

use self::menu::MenuItem;
use self::user::User;

pub mod reducer;
pub mod settings;
pub mod menu;
pub mod user;

pub type ReducerSender = UnboundedSender<ReducerEnum>;
pub type ReducerReceiver = UnboundedReceiver<ReducerEnum>;

#[enum_dispatch]
pub enum ReducerEnum {
    StartOrder,
    AddItem,
    Finalize,
    Cancel,
    Arrived,
    Help,
}

pub enum MachineState {
    Idle,
    TakeOrders,
    Ordered,
}

pub struct Store {
    rx: ReducerReceiver,
    tx: ReducerSender,
    reducers: Vec<ReducerEnum>,
    state: State
}

impl Store {
    pub fn get_sender(&self) -> ReducerSender {
        self.tx.clone()
    }
}

struct State {
    orders: HashMap<User, Vec<MenuItem>>,
    machine_state: MachineState
}

enum ReducerError {
    InvalidTransition {message: String}
}

#[async_trait]
#[enum_dispatch(ReducerEnum)]
trait Action {
    // passing in a mut state is fine since the reduce function is completely
    // internal and cloning and returning the same state is unnecessary overhead
    fn reduce(&self, state: &mut State) -> Result<(), ReducerError>;
}
