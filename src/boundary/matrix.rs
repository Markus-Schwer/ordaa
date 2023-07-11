use async_trait::async_trait;

use crate::control::store::{ActionSender, SharableState};

use super::Runnable;

pub struct MatrixBot {}

#[async_trait]
impl Runnable for MatrixBot {
    async fn run(&self, _sender: ActionSender, _state: SharableState) {
        println!("starting matrix bot");
    }
}
