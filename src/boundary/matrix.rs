use crate::control::ActionSender;
use crate::State;
use async_trait::async_trait;
use std::sync::Arc;
use tokio::sync::RwLock;

use super::RunnableBoundary;

pub struct MatrixBot {}

#[async_trait]
impl RunnableBoundary for MatrixBot {
    async fn run(&self, sender: ActionSender, state: Arc<RwLock<State>>) {
        println!("starting matrix bot");
    }
}
