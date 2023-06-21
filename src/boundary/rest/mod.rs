
use std::sync::Arc;

use crate::control::ActionSender;
use crate::State;

use super::RunnableBoundary;
use async_trait::async_trait;
use tokio::sync::RwLock;

pub mod routes;
mod handlers;

pub struct RestApi {}

#[async_trait]
impl RunnableBoundary for RestApi {
    async fn run(&self, sender: ActionSender, _state: Arc<RwLock<State>>) {
        println!("starting REST API on 127.0.0.1:8080");
        // Match any request and return hello world!
        let routes = routes::hello_routes(sender);

        warp::serve(routes).run(([127, 0, 0, 1], 8080)).await;
    }
}

