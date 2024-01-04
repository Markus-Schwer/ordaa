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

#[derive(Template)]
#[template(path = "menus.html")]
pub struct MenusTemplate;

#[derive(Template)]
#[template(path = "menu.html")]
pub struct MenuTemplate;

pub mod filters {
    use askama::Template;
    use warp::{Filter, reply::html};
    use super::{IndexTemplate, AdminTemplate, OrdersTemplate, OrderTemplate, MenusTemplate, MenuTemplate};

    pub fn all(
    ) -> impl Filter<Extract = (impl warp::Reply,), Error = warp::Rejection> + Clone {
        index().or(static_files()).or(orders()).or(order()).or(menus()).or(menu()).or(admin())
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

    fn menus() -> impl Filter<Extract = (impl warp::Reply,), Error = warp::Rejection> + Clone {
        warp::path!("menus")
            .and(warp::get())
            .and_then(|| async move {
                Ok::<warp::reply::Html<String>, warp::Rejection>(html(MenusTemplate {}.render().unwrap()))
            })
    }

    fn menu() -> impl Filter<Extract = (impl warp::Reply,), Error = warp::Rejection> + Clone {
        warp::path!("menu")
            .and(warp::get())
            .and_then(|| async move {
                Ok::<warp::reply::Html<String>, warp::Rejection>(html(MenuTemplate {}.render().unwrap()))
            })
    }
}
