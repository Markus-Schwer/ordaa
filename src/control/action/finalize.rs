use std::matches;

use crate::control::state::{MachineState, State};

pub struct Finalize {}

impl super::Action for Finalize {
    fn reduce(&self, mut state: State) -> Result<State, super::ReducerError> {
        if !matches!(state.machine_state, MachineState::TakeOrders) {
            return Err(super::ReducerError::InvalidTransition {
                message: "there is nothing to finalize right now".into(),
            });
        }
        if state.orders.is_empty() {
            return Err(super::ReducerError::InvalidState {
                message: "there are no orders, won't finalize an empty order".into(),
            });
        }
        state.machine_state = MachineState::Ordered;
        Ok(state)
    }
}
