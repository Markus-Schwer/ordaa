use std::convert::Infallible;
use warp;

use crate::control::ActionSender;

pub async fn hello(_sender: ActionSender) -> Result<impl warp::Reply, Infallible> {
    Ok(warp::reply::with_status("hello world", warp::http::StatusCode::OK))
}
