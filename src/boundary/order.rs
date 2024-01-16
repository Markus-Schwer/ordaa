use std::error::Error;

use actix_web::{get, post, web, Responder};
use crate::service::state::AppState;
use super::dto::NewOrderDto;

pub fn services_order(cfg: &mut web::ServiceConfig) {
    cfg.service(all_orders);
    cfg.service(get_order);
    cfg.service(create_order);
}

#[post("/order")]
async fn create_order(new_order: web::Json<NewOrderDto>, data: web::Data<AppState>) -> Result<impl Responder, Box<dyn Error>> {
    let order = data.db.create_order(new_order.into_inner());
    Ok(web::Json(order))
}

#[get("/order")]
async fn all_orders(data: web::Data<AppState>) -> Result<impl Responder, Box<dyn Error>> {
    Ok(web::Json(data.db.all_orders()))
}

#[get("/order/{order_id}")]
async fn get_order(path: web::Path<(i32,)>, data: web::Data<AppState>) -> Result<impl Responder, Box<dyn Error>> {
    Ok(web::Json(data.db.all_order_items(path.0)))
}
