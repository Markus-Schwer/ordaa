use tokio::sync::mpsc::{UnboundedSender, UnboundedReceiver, self};

use super::state::{Transition, State};

pub struct StateMachine {
    transition_rx: UnboundedReceiver<Transition>,
    transition_tx: UnboundedSender<Transition>,
    current_state: State
}

#[derive(Clone)]
pub struct ReadOnly {
    transition_tx: UnboundedSender<Transition>,
}

impl StateMachine {
    pub fn new() -> Self {
        // if we want backpressure, we should use channel instead of unbounded_channel
        let (tx, rx) = mpsc::unbounded_channel();
        StateMachine { transition_rx: rx, transition_tx: tx, current_state: State::Idle }
    }
    pub fn get_clone(&self) -> ReadOnly {
        ReadOnly { transition_tx: self.transition_tx.clone() }
    }
    pub fn start(&mut self) {
        loop {
            if let msg = Some(self.transition_rx.recv()) {
                // do transition
            } else {
                // handle error
            }
        }
    }
}
