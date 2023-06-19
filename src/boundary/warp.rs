#![deny(warnings)]
use warp::Filter;
use crate::control::machine::ReadOnly;

use super::RunnableBoundary;
use async_trait::async_trait;

pub struct WarpBoundary {}

#[async_trait]
impl RunnableBoundary for WarpBoundary {
    async fn run(&self, _state_machine: ReadOnly) {
        println!("starting REST API");
        // Match any request and return hello world!
        let routes = warp::any().map(|| "Hello, World!");

        warp::serve(routes).run(([127, 0, 0, 1], 8080)).await;
    }
}

