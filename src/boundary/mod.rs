pub mod matrix;
pub mod rest;
use async_trait::async_trait;
use enum_dispatch::enum_dispatch;

use crate::control::machine::ReadOnly;
use crate::boundary::rest::RestApi;
use crate::boundary::matrix::MatrixBot;

#[enum_dispatch]
pub enum BoundaryEnum {
    RestApi,
    MatrixBot
}

#[async_trait]
#[enum_dispatch(BoundaryEnum)]
pub trait RunnableBoundary {
    async fn run(&self, state_machine: ReadOnly);
}
