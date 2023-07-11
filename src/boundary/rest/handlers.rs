use std::convert::Infallible;
use warp;

use crate::{
    boundary::SharableState,
    control::{
        action::{help::Help, start_order::StartOrder, Action},
        store::ActionSender,
    },
};

pub async fn help(
    sender: ActionSender,
    state: SharableState,
) -> Result<impl warp::Reply, Infallible> {
    if let Err(err) = sender.send(Action::from(Help {})) {
        Ok(warp::reply::with_status(
            err.to_string(),
            warp::http::StatusCode::INTERNAL_SERVER_ERROR,
        ))
    } else {
        Ok(warp::reply::with_status(
            state.read().await.to_string(),
            warp::http::StatusCode::OK,
        ))
    }
}

pub async fn start_order(sender: ActionSender) -> Result<impl warp::Reply, Infallible> {
    if let Err(err) = sender.send(Action::from(StartOrder {})) {
        Ok(warp::reply::with_status(
            err.to_string(),
            warp::http::StatusCode::INTERNAL_SERVER_ERROR,
        ))
    } else {
        Ok(warp::reply::with_status(
            "OK".to_string(),
            warp::http::StatusCode::OK,
        ))
    }
}
