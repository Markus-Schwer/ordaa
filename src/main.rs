use std::sync::Arc;

use tokio::sync::RwLock;
use tokio::task::JoinSet;

use crate::boundary::matrix::MatrixBot;
use crate::boundary::rest::RestApi;
use crate::boundary::{BoundaryEnum, RunnableBoundary};
use crate::control::{state::State, store::Store};

mod boundary;
mod control;

#[tokio::main]
async fn main() {
    println!("Starting .inder server");

    let state = Arc::new(RwLock::new(State::new()));
    let mut store = Store::new(state.clone());

    // setup boundaries
    let boundaries: Vec<BoundaryEnum> = vec![RestApi {}.into(), MatrixBot {}.into()];

    let mut join_set = JoinSet::new();
    for boundary in boundaries {
        let sender = store.get_sender();
        let s = state.clone();
        join_set.spawn(async move {
            boundary.run(sender, s).await;
        });
    }

    // finally start the state machine
    join_set.spawn(async move {
        store.listen().await;
    });

    while let Some(res) = join_set.join_next().await {
        res.unwrap();
    }
}
