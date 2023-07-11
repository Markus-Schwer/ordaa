use async_trait::async_trait;

use crate::control::store::ActionSender;
use super::{Runnable, SharableState};

mod handlers;
pub mod routes;

pub struct RestApi {}

#[async_trait]
impl Runnable for RestApi {
    async fn run(&self, sender: ActionSender, state: SharableState) {
        println!("starting REST API on 127.0.0.1:8080");
        let routes = routes::rest_routes(sender, state);
        warp::serve(routes).run(([127, 0, 0, 1], 8080)).await;
    }
}
