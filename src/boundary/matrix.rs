use async_trait::async_trait;
use crate::control::machine::ReadOnly;
use super::RunnableBoundary;

pub struct MatrixBot {
}

#[async_trait]
impl RunnableBoundary for MatrixBot {
    async fn run(&self, _state_machine: ReadOnly) {
        println!("starting matrix bot");
    }
}
