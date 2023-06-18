use std::sync::mpsc::{Sender, Receiver, self};

use super::state::{Transition, State};

pub struct StateMachine {
    transition_rx: Receiver<Transition>,
    transition_tx: Sender<Transition>,
    current_state: State
}

#[derive(Clone)]
pub struct ReadOnly {
    transition_tx: Sender<Transition>,
}

impl StateMachine {
    pub fn new() -> Self {
        let (tx, rx) = mpsc::channel();
        StateMachine { transition_rx: rx, transition_tx: tx, current_state: State::Idle }
    }
    pub fn get_clone(&self) -> ReadOnly {
        ReadOnly { transition_tx: self.transition_tx.clone() }
    }
    pub fn start(&self) {
        loop {
            if let msg = Some(self.transition_rx.recv()) {
                // do transition
            } else {
                // handle error
            }
        }
    }
}
