pub mod matrix;
pub mod rest;

use async_trait::async_trait;
use enum_dispatch::enum_dispatch;

use crate::boundary::matrix::MatrixBot;
use crate::boundary::rest::RestApi;
use crate::control::store::{ActionSender, SharableState};

#[enum_dispatch]
pub enum Boundary {
    RestApi,
    MatrixBot,
}

#[async_trait]
#[enum_dispatch(Boundary)]
pub trait Runnable {
    async fn run(&self, sender: ActionSender, state: SharableState);
}
