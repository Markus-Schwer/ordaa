use std::matches;

use crate::control::{Action, MachineState, ReducerError};

pub struct StartOrder {}

impl Action for StartOrder {
    fn reduce(
        &self,
        mut state: crate::control::State,
    ) -> Result<crate::control::State, crate::control::ReducerError> {
        if matches!(state.machine_state, MachineState::Idle) {
            state.machine_state = MachineState::TakeOrders;
        } else {
            return Err(ReducerError::InvalidTransition {
                message: "cannot start ordering right now".into(),
            });
        }
        Ok(state)
    }
}
