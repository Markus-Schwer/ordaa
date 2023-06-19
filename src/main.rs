mod boundary;
mod control;
use crate::control::machine::StateMachine;
use crate::boundary::{BoundaryEnum, RunnableBoundary};
use crate::boundary::rest::RestApi;
use crate::boundary::matrix::MatrixBot;
use tokio::task::JoinSet;

#[tokio::main]
async fn main() {
    println!("Starting .inder server");

    let mut sm = StateMachine::new();

    // setup boundaries
    let boundaries: Vec<BoundaryEnum> = vec![RestApi {}.into(), MatrixBot {}.into()];

    let mut join_set = JoinSet::new();
    for boundary in boundaries {
        let cloned_sm = sm.get_clone();
        join_set.spawn(async move {
            boundary.run(cloned_sm).await;
        });
    }

    // finally start the state machine
    join_set.spawn(async move {
        sm.start();
    });

    while let Some(res) = join_set.join_next().await {
        res.unwrap();
    }
}
