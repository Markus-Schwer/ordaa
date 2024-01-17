use std::error::Error;

use diesel::prelude::*;
use actix_web::{get, post, web, Responder};
use crate::{service::state::AppState, boundary::dto::NewOrderItemDto};
use super::dto::NewOrderDto;

pub fn services_order(cfg: &mut web::ServiceConfig) {
    cfg.service(all_orders);
    cfg.service(get_order);
    cfg.service(get_order_items);
    cfg.service(create_order);
    cfg.service(create_order_item);
}

#[post("/order")]
async fn create_order(new_order: web::Json<NewOrderDto>, data: web::Data<AppState>) -> Result<impl Responder, Box<dyn Error>> {
    data.db.get_conn()?.transaction(|conn| {
        let order = data.db.create_order(conn, new_order.into_inner())?;
        Ok(web::Json(order))
    })
}

#[get("/order")]
async fn all_orders(data: web::Data<AppState>) -> Result<impl Responder, Box<dyn Error>> {
    data.db.get_conn()?.transaction(|conn| {
        Ok(web::Json(data.db.all_orders(conn)?))
    })
}

#[get("/order/{order_id}")]
async fn get_order(path: web::Path<(i32,)>, data: web::Data<AppState>) -> Result<impl Responder, Box<dyn Error>> {
    data.db.get_conn()?.transaction(|conn| {
        Ok(web::Json(data.db.get_order_by_id(conn, path.0)?))
    })
}

#[get("/order/{order_id}/order-item")]
async fn get_order_items(path: web::Path<(i32,)>, data: web::Data<AppState>) -> Result<impl Responder, Box<dyn Error>> {
    data.db.get_conn()?.transaction(|conn| {
        Ok(web::Json(data.db.all_order_items(conn, path.0)?))
    })
}

#[post("/order/{order_id}/order-item")]
async fn create_order_item(new_order_item: web::Json<NewOrderItemDto>, data: web::Data<AppState>) -> Result<impl Responder, Box<dyn Error>> {
    data.db.get_conn()?.transaction(|conn| {
        let order_item = data.db.create_order_item(conn, new_order_item.into_inner())?;
        Ok(web::Json(order_item))
    })
}
