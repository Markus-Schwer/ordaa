use std::sync::Arc;
use async_trait::async_trait;
use tokio::sync::RwLock;
use crate::control::ActionSender;
use crate::State;

use super::RunnableBoundary;

pub struct MatrixBot {
}

#[async_trait]
impl RunnableBoundary for MatrixBot {
    async fn run(&self, sender: ActionSender, state: Arc<RwLock<State>>) {
        println!("starting matrix bot");
    }
}
