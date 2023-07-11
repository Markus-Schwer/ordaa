use std::{collections::HashMap, format};

#[derive(Clone, Debug)]
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

impl ToString for State {
    fn to_string(&self) -> String {
        let mut res = self
            .orders
            .iter()
            .map(|order| format!("user: {}; orders: {:?}", order.0.to_string(), order.1))
            .collect::<Vec<String>>()
            .join("\n");
        res.push_str(format!("\ncurrent state: {:?}", self.machine_state).as_str());
        return res;
    }
}
