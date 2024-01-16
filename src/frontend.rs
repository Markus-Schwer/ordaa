use std::error::Error;

use actix_web::{web, Responder, get};
use actix_files::Files;
use askama::Template;

use crate::{boundary::dto::MenuDto, service::state::AppState};

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
    pub menus: Vec<MenuDto>
}

#[derive(Template)]
#[template(path = "menu.html")]
pub struct MenuTemplate {
    pub menu: MenuDto,
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
    Ok(MenusTemplate { menus: data.db.all_menus() })
}

#[get("/menu/{menu_id}")]
async fn get_menu(path: web::Path<(i32,)>, data: web::Data<AppState>) -> Result<impl Responder, Box<dyn Error>> {
    Ok(MenuTemplate { menu: data.db.get_menu_by_id(path.0) })
}

#[get("/orders")]
async fn get_orders(data: web::Data<AppState>) -> Result<impl Responder, Box<dyn Error>> {
    Ok(OrdersTemplate { })
}

#[get("/order/{order_id}")]
async fn get_order(path: web::Path<(i32,)>, data: web::Data<AppState>) -> Result<impl Responder, Box<dyn Error>> {
    Ok(OrderTemplate { })
}
