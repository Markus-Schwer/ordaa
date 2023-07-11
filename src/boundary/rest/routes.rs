use warp::{self, Filter, Rejection, Reply};

use crate::{
    boundary::{rest::handlers, SharableState},
    control::store::ActionSender,
};

pub fn rest_routes(
    sender: ActionSender,
    state: SharableState,
) -> impl Filter<Extract = (impl Reply,), Error = Rejection> + Clone {
    warp::any()
        .and(help(sender.clone(), state.clone()))
        .or(start_order(sender.clone()))
}

fn start_order(
    sender: ActionSender,
) -> impl Filter<Extract = (impl Reply,), Error = Rejection> + Clone {
    warp::path("start-order")
        .and(warp::get())
        .and(warp::any().map(move || sender.clone()))
        .and_then(handlers::start_order)
}

fn help(
    sender: ActionSender,
    state: SharableState,
) -> impl Filter<Extract = (impl Reply,), Error = Rejection> + Clone {
    warp::path("help")
        .and(warp::get())
        .and(warp::any().map(move || sender.clone()))
        .and(warp::any().map(move || state.clone()))
        .and_then(handlers::help)
}
