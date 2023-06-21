use std::matches;

use crate::control::{Action, MachineState};

pub struct Finalize {}

impl Action for Finalize {
    fn reduce(
        &self,
        mut state: crate::control::State,
    ) -> Result<crate::control::State, crate::control::ReducerError> {
        if !matches!(state.machine_state, MachineState::TakeOrders) {
            return Err(crate::control::ReducerError::InvalidTransition {
                message: "there is nothing to finalize right now".into(),
            });
        }
        if state.orders.is_empty() {
            return Err(crate::control::ReducerError::InvalidState {
                message: "there are no orders, won't finalize an empty order".into(),
            });
        }
        state.machine_state = MachineState::Ordered;
        Ok(state)
    }
}
