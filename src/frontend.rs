use std::error::Error;

use actix_web::{web, Responder, get};
use actix_files::Files;
use askama::Template;
use itertools::Itertools;
use sqlx::Connection;
use uuid::Uuid;

use crate::{service::state::AppState, entity::models::{User, MenuWithItems, OrderWithItems, OrderItemWithJoins}};

#[derive(Template)]
#[template(path = "index.html")]
pub struct IndexTemplate;

#[derive(Template)]
#[template(path = "orders.html")]
pub struct OrdersTemplate {
    pub orders: Vec<OrderWithItems>
}

#[derive(Template)]
#[template(path = "order.html")]
pub struct OrderTemplate {
    pub order: OrderWithItems,
    pub price_total: i32,
    pub grouped_items: Vec<(User, i32, Vec<OrderItemWithJoins>)>
}

#[derive(Template)]
#[template(path = "admin.html")]
pub struct AdminTemplate;

#[derive(Template)]
#[template(path = "menus.html")]
pub struct MenusTemplate {
    pub menus: Vec<MenuWithItems>
}

#[derive(Template)]
#[template(path = "menu.html")]
pub struct MenuTemplate {
    pub menu: MenuWithItems,
}

pub fn services_frontend(cfg: &mut web::ServiceConfig) {
    cfg.service(web::resource("/").to(|| async { IndexTemplate {} }));
    cfg.service(get_orders);
    cfg.service(get_order);
    cfg.service(get_menus);
    cfg.service(get_menu);
    cfg.service(web::resource("/admin").to(|| async { AdminTemplate {} }));
    cfg.service(Files::new("/static", "./static").prefer_utf8(true));
}

#[get("/menus")]
async fn get_menus(data: web::Data<AppState>) -> Result<impl Responder, Box<dyn Error>> {
    let mut conn = data.db.get_conn().await?;
    let mut tx = conn.begin().await?;

    let menus = data.db.all_menus(&mut tx).await?;

    tx.rollback().await?;
    conn.close().await?;
    Ok(MenusTemplate { menus })
}

#[get("/menu/{menu_id}")]
async fn get_menu(path: web::Path<(Uuid,)>, data: web::Data<AppState>) -> Result<impl Responder, Box<dyn Error>> {
    let mut conn = data.db.get_conn().await?;
    let mut tx = conn.begin().await?;

    let menu = data.db.get_menu_by_uuid(&mut tx, path.0).await?;

    tx.rollback().await?;
    conn.close().await?;
    Ok(MenuTemplate { menu })
}

#[get("/orders")]
async fn get_orders(data: web::Data<AppState>) -> Result<impl Responder, Box<dyn Error>> {
    let mut conn = data.db.get_conn().await?;
    let mut tx = conn.begin().await?;

    let orders = data.db.all_orders(&mut tx).await?;

    tx.rollback().await?;
    conn.close().await?;
    Ok(OrdersTemplate { orders })
}

#[get("/order/{order_id}")]
async fn get_order(path: web::Path<(Uuid,)>, data: web::Data<AppState>) -> Result<impl Responder, Box<dyn Error>> {
    let mut conn = data.db.get_conn().await?;
    let mut tx = conn.begin().await?;

    let order = data.db.get_order_by_uuid(&mut tx, path.0).await?;
    let price_total: i32 = order.items.iter().map(|oi| oi.price).sum();
    let grouped_items: Vec<(User, i32, Vec<OrderItemWithJoins>)> = order.items.iter().group_by(|elt: &&OrderItemWithJoins| elt.order_user.clone()).into_iter()
        .map(|(ge0, group)| {
            let items: Vec<OrderItemWithJoins> = group.cloned().collect();
            let group_total = items.iter().map(|oi| oi.price).sum();
            (ge0, group_total, items)
        })
        .collect();

    tx.rollback().await?;
    conn.close().await?;
    Ok(OrderTemplate { order, price_total, grouped_items })
}
