use std::matches;

use crate::control::{Action, MachineState, ReducerError};

pub struct Arrived {}

impl Action for Arrived {
    fn reduce(
        &self,
        mut state: crate::control::State,
    ) -> Result<crate::control::State, crate::control::ReducerError> {
        if matches!(state.machine_state, MachineState::Ordered) {
            state.machine_state = MachineState::Idle;
        } else {
            return Err(ReducerError::InvalidTransition {
                message: "what just arrived? nothing I ordered I guess".into(),
            });
        }
        Ok(state)
    }
}
