use core::panic;
use enum_dispatch::enum_dispatch;
use std::collections::HashMap;
use std::vec;
use tokio::sync::mpsc::{self, UnboundedReceiver, UnboundedSender};

mod action;
use action::add_item::AddItem;
use action::arrived::Arrived;
use action::cancel::Cancel;
use action::finalize::Finalize;
use action::help::Help;
use action::start_order::StartOrder;

use self::menu::MenuItem;
use self::user::User;

pub mod menu;
pub mod settings;
pub mod user;

pub type ActionSender = UnboundedSender<ActionEnum>;
pub type ActionReceiver = UnboundedReceiver<ActionEnum>;
pub type EffectFn = fn(state: &State);

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

#[derive(Clone)]
pub enum MachineState {
    Idle,
    TakeOrders,
    Ordered,
}

pub struct Store {
    rx: ActionReceiver,
    tx: ActionSender,
    effects: HashMap<ActionEnum, Vec<EffectFn>>,
    state: State,
}

impl Store {
    pub fn new() -> Self {
        let (tx, rx) = mpsc::unbounded_channel();
        Self {
            rx,
            tx,
            effects: HashMap::new(),
            state: State::new(),
        }
    }
    pub fn get_sender(&self) -> ActionSender {
        self.tx.clone()
    }
    pub fn register_effect(&mut self, effect: EffectFn, for_action: ActionEnum) {
        if let Some(effects) = self.effects.get(&for_action) {
            effects.push(effect);
        } else {
            self.effects.insert(for_action, vec![effect]);
        }
    }
    pub async fn listen(&mut self) {
        loop {
            if let Some(action) = self.rx.recv().await {
                match action.reduce(self.state.clone()) {
                    Ok(new_state) => {
                        self.state = new_state;
                        if let Some(effects) = self.effects.get(&action) {
                            for effect in effects {
                                effect(&self.state);
                            }
                        }
                    }
                    Err(err) => panic!("{:?}", err),
                }
            } else {
                panic!("received empty action signal");
            }
        }
    }
}

#[derive(Clone)]
pub struct State {
    orders: HashMap<User, Vec<MenuItem>>,
    machine_state: MachineState,
}

impl State {
    pub fn new() -> Self {
        Self {
            orders: HashMap::new(),
            machine_state: MachineState::Idle,
        }
    }
}

#[derive(Debug)]
enum ReducerError {
    InvalidTransition { message: String },
    InvalidState { message: String },
}

#[enum_dispatch(ActionEnum)]
trait Action {
    // it feels bad to copy the state every call, but passing it mutably is
    // just as bad.
    // TODO: optimize this if possible
    fn reduce(&self, state: State) -> Result<State, ReducerError>;
}
