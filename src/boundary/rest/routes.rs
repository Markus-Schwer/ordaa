use warp::{self, Filter, Reply, Rejection};
use crate::control::machine::ReadOnly;
use crate::boundary::rest::handlers;
use std::convert::Infallible;

fn with_state_machine(sm: ReadOnly) -> impl Filter<Extract = (ReadOnly,), Error = Infallible> + Clone {
    warp::any().map(move || sm.clone())
}

pub fn hello_routes(sm: ReadOnly) -> impl Filter<Extract = (impl Reply,), Error = Rejection> + Clone {
    hello(sm.clone())
}

// GET /hello
fn hello(sm: ReadOnly) -> impl Filter<Extract = (impl Reply,), Error = Rejection> + Clone {
    warp::path("hello").and(warp::get()).and(with_state_machine(sm)).and_then(handlers::hello)
}
