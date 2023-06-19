pub mod matrix;
pub mod warp;
use async_trait::async_trait;
use enum_dispatch::enum_dispatch;

use crate::control::machine::ReadOnly;
use crate::boundary::warp::WarpBoundary;
use crate::boundary::matrix::MatrixBot;

#[enum_dispatch]
pub enum BoundaryEnum {
    WarpBoundary,
    MatrixBot
}

#[async_trait]
#[enum_dispatch(BoundaryEnum)]
pub trait RunnableBoundary {
    async fn run(&self, state_machine: ReadOnly);
}
