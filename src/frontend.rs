use askama::Template;

#[derive(Template)]
#[template(path = "index.html")]
pub struct IndexTemplate;

#[derive(Template)]
#[template(path = "orders.html")]
pub struct OrdersTemplate;

#[derive(Template)]
#[template(path = "order.html")]
pub struct OrderTemplate;

#[derive(Template)]
#[template(path = "admin.html")]
pub struct AdminTemplate;

pub mod filters {
    use askama::Template;
    use sqlx::SqlitePool;
    use warp::{Filter, reply::html};
    use super::{IndexTemplate, AdminTemplate, OrdersTemplate, OrderTemplate};

    pub fn all(
        _pool: &SqlitePool,
    ) -> impl Filter<Extract = (impl warp::Reply,), Error = warp::Rejection> + Clone {
        index().or(static_files()).or(orders()).or(order()).or(admin())
    }

    fn static_files() -> impl Filter<Extract = (impl warp::Reply,), Error = warp::Rejection> + Clone {
        warp::path("static").and(warp::fs::dir("static"))
    }

    fn index() -> impl Filter<Extract = (impl warp::Reply,), Error = warp::Rejection> + Clone {
        warp::path::end()
            .and(warp::get())
            .and_then(|| async move {
                Ok::<warp::reply::Html<String>, warp::Rejection>(html(IndexTemplate {}.render().unwrap()))
            })
    }

    fn orders() -> impl Filter<Extract = (impl warp::Reply,), Error = warp::Rejection> + Clone {
        warp::path!("orders")
            .and(warp::get())
            .and_then(|| async move {
                Ok::<warp::reply::Html<String>, warp::Rejection>(html(OrdersTemplate {}.render().unwrap()))
            })
    }

    fn order() -> impl Filter<Extract = (impl warp::Reply,), Error = warp::Rejection> + Clone {
        warp::path!("order")
            .and(warp::get())
            .and_then(|| async move {
                Ok::<warp::reply::Html<String>, warp::Rejection>(html(OrderTemplate {}.render().unwrap()))
            })
    }

    fn admin() -> impl Filter<Extract = (impl warp::Reply,), Error = warp::Rejection> + Clone {
        warp::path!("admin")
            .and(warp::get())
            .and_then(|| async move {
                Ok::<warp::reply::Html<String>, warp::Rejection>(html(AdminTemplate {}.render().unwrap()))
            })
    }
}
