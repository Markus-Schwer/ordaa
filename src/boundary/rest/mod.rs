use crate::control::machine::ReadOnly;

use super::RunnableBoundary;
use async_trait::async_trait;

pub mod routes;
mod handlers;

pub struct RestApi {}

#[async_trait]
impl RunnableBoundary for RestApi {
    async fn run(&self, state_machine: ReadOnly) {
        println!("starting REST API on 127.0.0.1:8080");
        // Match any request and return hello world!
        let routes = routes::hello_routes(state_machine);

        warp::serve(routes).run(([127, 0, 0, 1], 8080)).await;
    }
}

