use std::collections::HashMap;

#[derive(Clone)]
pub enum MachineState {
    Idle,
    TakeOrders,
    Ordered,
}

#[derive(Clone)]
pub struct State {
    pub orders: HashMap<super::user::User, Vec<super::menu::MenuItem>>,
    pub machine_state: MachineState,
}

impl State {
    pub fn new() -> Self {
        Self {
            orders: HashMap::new(),
            machine_state: MachineState::Idle,
        }
    }
}
