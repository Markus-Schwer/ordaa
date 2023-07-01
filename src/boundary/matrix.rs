use std::sync::Arc;

use async_trait::async_trait;
use tokio::sync::RwLock;

use crate::control::{store::ActionSender, state::State};

use super::RunnableBoundary;

pub struct MatrixBot {}

#[async_trait]
impl RunnableBoundary for MatrixBot {
    async fn run(&self, _sender: ActionSender, _state: Arc<RwLock<State>>) {
        println!("starting matrix bot");
    }
}
