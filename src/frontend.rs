use std::error::Error;

use itertools::Itertools;
use diesel::prelude::*;
use actix_web::{web, Responder, get};
use actix_files::Files;
use askama::Template;

use crate::{boundary::dto::{OrderDto, UserDto, OrderItemDto, MenuWithItemsDto}, service::state::AppState};

#[derive(Template)]
#[template(path = "index.html")]
pub struct IndexTemplate;

#[derive(Template)]
#[template(path = "orders.html")]
pub struct OrdersTemplate {
    pub orders: Vec<OrderDto>
}

#[derive(Template)]
#[template(path = "order.html")]
pub struct OrderTemplate {
    pub order: OrderDto,
    pub price_total: i32,
    pub grouped_items: Vec<(UserDto, i32, Vec<OrderItemDto>)>
}

#[derive(Template)]
#[template(path = "admin.html")]
pub struct AdminTemplate;

#[derive(Template)]
#[template(path = "menus.html")]
pub struct MenusTemplate {
    pub menus: Vec<MenuWithItemsDto>
}

#[derive(Template)]
#[template(path = "menu.html")]
pub struct MenuTemplate {
    pub menu: MenuWithItemsDto,
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
    data.db.get_conn()?.transaction(|conn| {
        Ok(MenusTemplate { menus: data.db.all_menus(conn)? })
    })
}

#[get("/menu/{menu_id}")]
async fn get_menu(path: web::Path<(i32,)>, data: web::Data<AppState>) -> Result<impl Responder, Box<dyn Error>> {
    data.db.get_conn()?.transaction(|conn| {
        Ok(MenuTemplate { menu: data.db.get_menu_by_id(conn, path.0)? })
    })
}

#[get("/orders")]
async fn get_orders(data: web::Data<AppState>) -> Result<impl Responder, Box<dyn Error>> {
    data.db.get_conn()?.transaction(|conn| {
        Ok(OrdersTemplate { orders: data.db.all_orders(conn)? })
    })
}

#[get("/order/{order_id}")]
async fn get_order(path: web::Path<(i32,)>, data: web::Data<AppState>) -> Result<impl Responder, Box<dyn Error>> {
    data.db.get_conn()?.transaction(|conn| {
        let order = data.db.get_order_by_id(conn, path.0)?;
        let price_total: i32 = order.items.iter().map(|oi| oi.price).sum();
        let grouped_items: Vec<(UserDto, i32, Vec<OrderItemDto>)> = order.items.iter().group_by(|elt: &&OrderItemDto| elt.user.clone()).into_iter()
            .map(|(ge0, group)| {
                let items: Vec<OrderItemDto> = group.cloned().collect();
                let group_total = items.iter().map(|oi| oi.price).sum();
                (ge0, group_total, items)
            })
            .collect();
        Ok(OrderTemplate { order, price_total, grouped_items })
    })
}
