use std::{matches, vec};

use crate::control::{
    menu::MenuItem,
    state::{MachineState, State},
    user::User,
};

pub struct AddItem {
    user: User,
    menu_item: MenuItem,
}

impl super::Action for AddItem {
    fn reduce(&self, mut state: State) -> Result<State, super::ReducerError> {
        if !matches!(state.machine_state, MachineState::TakeOrders) {
            return Err(super::ReducerError::InvalidTransition {
                message: "cannot place orders right now".into(),
            });
        }
        if let Some(user_orders) = state.orders.get_mut(&self.user) {
            user_orders.push(self.menu_item.clone())
        } else {
            state
                .orders
                .insert(self.user.clone(), vec![self.menu_item.clone()]);
        }
        Ok(state)
    }
}
