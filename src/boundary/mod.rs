pub mod matrix;
pub mod rest;
use std::sync::Arc;

use crate::State;

use async_trait::async_trait;
use enum_dispatch::enum_dispatch;
use tokio::sync::RwLock;

use crate::boundary::matrix::MatrixBot;
use crate::boundary::rest::RestApi;
use crate::control::ActionSender;

#[enum_dispatch]
pub enum BoundaryEnum {
    RestApi,
    MatrixBot,
}

#[async_trait]
#[enum_dispatch(BoundaryEnum)]
pub trait RunnableBoundary {
    async fn run(&self, sender: ActionSender, state: Arc<RwLock<State>>);
}
