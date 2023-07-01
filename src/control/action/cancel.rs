use std::matches;

use crate::control::{user::User, Action, MachineState, ReducerError};

pub struct Cancel {
    // if user is provided cancel orders of that user, otherwise cancel entire order
    user: Option<User>,
}

impl Action for Cancel {
    fn reduce(
        &self,
        mut state: crate::control::State,
    ) -> Result<crate::control::State, crate::control::ReducerError> {
        if !matches!(state.machine_state, MachineState::TakeOrders) {
            return Err(ReducerError::InvalidTransition {
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
