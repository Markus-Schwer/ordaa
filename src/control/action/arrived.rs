use std::matches;

use crate::control::state::{MachineState, State};

pub struct Arrived {}

impl super::Reducer for Arrived {
    fn reduce(&self, mut state: State) -> Result<State, super::ReducerError> {
        if matches!(state.machine_state, MachineState::Ordered) {
            state.machine_state = MachineState::Idle;
        } else {
            return Err(super::ReducerError::InvalidTransition {
                message: "what just arrived? nothing I ordered I guess".into(),
            });
        }
        Ok(state)
    }
}
