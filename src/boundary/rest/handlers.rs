use std::convert::Infallible;
use warp;

use crate::control::machine::ReadOnly;


pub async fn hello(sm: ReadOnly) -> Result<impl warp::Reply, Infallible> {
    Ok(warp::reply::with_status("hello world", warp::http::StatusCode::OK))
}
