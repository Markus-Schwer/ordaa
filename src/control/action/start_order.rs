use std::matches;

use crate::control::state::{MachineState, State};

pub struct StartOrder {}

impl super::Action for StartOrder {
    fn reduce(&self, mut state: State) -> Result<State, super::ReducerError> {
        if matches!(state.machine_state, MachineState::Idle) {
            state.machine_state = MachineState::TakeOrders;
        } else {
            return Err(super::ReducerError::InvalidTransition {
                message: "cannot start ordering right now".into(),
            });
        }
        Ok(state)
    }
}
