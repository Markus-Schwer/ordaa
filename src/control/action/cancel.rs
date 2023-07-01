use std::matches;

use crate::control::state::{MachineState, State};

pub struct Cancel {
    // if user is provided cancel orders of that user, otherwise cancel entire order
    user: Option<crate::control::user::User>,
}

impl super::Action for Cancel {
    fn reduce(&self, mut state: State) -> Result<State, super::ReducerError> {
        if !matches!(state.machine_state, MachineState::TakeOrders) {
            return Err(super::ReducerError::InvalidTransition {
                message: "there is nothing to cancel or it's already too late".into(),
            });
        }
        if let Some(user) = &self.user {
            state.orders.remove(&user);
        } else {
            state.orders.clear();
            state.machine_state = MachineState::Idle;
        }
        Ok(state)
    }
}
