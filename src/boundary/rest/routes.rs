use crate::boundary::rest::handlers;
use crate::control::ActionSender;
use std::convert::Infallible;
use warp::{self, Filter, Rejection, Reply};

fn with_sender(
    sender: ActionSender,
) -> impl Filter<Extract = (ActionSender,), Error = Infallible> + Clone {
    warp::any().map(move || sender.clone())
}

pub fn hello_routes(
    sender: ActionSender,
) -> impl Filter<Extract = (impl Reply,), Error = Rejection> + Clone {
    hello(sender.clone())
}

// GET /hello
fn hello(sender: ActionSender) -> impl Filter<Extract = (impl Reply,), Error = Rejection> + Clone {
    warp::path("hello")
        .and(warp::get())
        .and(with_sender(sender))
        .and_then(handlers::hello)
}
