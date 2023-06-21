use std::{vec, matches};

use crate::control::{menu::MenuItem, user::User, Action, MachineState, ReducerError, State};

pub struct AddItem {
    user: User,
    menu_item: MenuItem,
}

impl Action for AddItem {
    fn reduce(
        &self,
        mut state: State,
    ) -> Result<State, crate::control::ReducerError> {
        if !matches!(state.machine_state, MachineState::TakeOrders) {
            return Err(ReducerError::InvalidTransition { message: "cannot place orders right now".into() });
        }
        if let Some(user_orders) = state.orders.get(&self.user) {
            user_orders.push(self.menu_item)
        } else {
            state.orders.insert(self.user, vec![self.menu_item]);
        }
        Ok(state)
    }
}
