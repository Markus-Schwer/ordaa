use askama::Template;

use crate::menu::Menu;

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
pub struct MenusTemplate {
    pub menus: Vec<Menu>
}

pub struct Item {
    pub id: String,
    pub name: String,
    pub price: String,
}

#[derive(Template)]
#[template(path = "menu.html")]
pub struct MenuTemplate {
    pub name: String,
    pub items: Vec<Item>
}

pub mod filters {
    use askama::Template;
    use warp::{Filter, reply::html};

    use crate::{menu::Menu, db::Db, search::SearchContextReader, filters::{with_db, with_searcher_ctx}};

    use super::{IndexTemplate, AdminTemplate, OrdersTemplate, OrderTemplate, MenusTemplate, MenuTemplate, Item};

    pub fn all(
        db: Db,
        ctx: SearchContextReader,
    ) -> impl Filter<Extract = (impl warp::Reply,), Error = warp::Rejection> + Clone {
        index().or(static_files()).or(orders()).or(order()).or(menus(db.clone(), ctx.clone())).or(menu(db.clone(), ctx.clone())).or(admin())
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

    fn menus(
        db: Db,
        ctx: SearchContextReader,
    ) -> impl Filter<Extract = (impl warp::Reply,), Error = warp::Rejection> + Clone {
        warp::path!("menus")
            .and(warp::get())
            .and(with_db(db.clone()))
            .and(with_searcher_ctx(ctx))
            .and_then(|db: Db, ctx: SearchContextReader| async move {
                Ok::<warp::reply::Html<String>, warp::Rejection>(html(MenusTemplate {
                    menus: vec![Menu { name: "Sangam".into(), items: vec![]}]
                }.render().unwrap()))
            })
    }

    fn menu(
        db: Db,
        ctx: SearchContextReader,
    ) -> impl Filter<Extract = (impl warp::Reply,), Error = warp::Rejection> + Clone {
        warp::path!("menu")
            .and(warp::get())
            .and(with_db(db.clone()))
            .and(with_searcher_ctx(ctx))
            .and_then(|db: Db, ctx: SearchContextReader| async move {
                Ok::<warp::reply::Html<String>, warp::Rejection>(html(MenuTemplate {
                    name: "Sangam".into(),
                    items: vec![
                        Item { id: "42".into(), name: "Chicken Tikka".into(), price: format!("{},{}€", 15, 19) }
                    ]
                }.render().unwrap()))
            })
    }
}
