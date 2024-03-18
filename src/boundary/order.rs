use std::error::Error;
use uuid::Uuid;

use actix_web::{get, post, web, Responder};
use sqlx::Connection;
use crate::{service::state::AppState, entity::models::{NewOrder, NewOrderItem}};

pub fn services_order(cfg: &mut web::ServiceConfig) {
    cfg.service(all_orders);
    cfg.service(get_order);
    cfg.service(get_order_items);
    cfg.service(create_order);
    cfg.service(create_order_item);
}

#[post("/order")]
async fn create_order(new_order: web::Json<NewOrder>, data: web::Data<AppState>) -> Result<impl Responder, Box<dyn Error>> {
    let mut conn = data.db.get_conn().await?;
    let mut tx = conn.begin().await?;

    let order = data.db.create_order(&mut tx, new_order.into_inner()).await?;

    tx.commit().await?;
    conn.close().await?;
    Ok(web::Json(order))
}

#[get("/order")]
async fn all_orders(data: web::Data<AppState>) -> Result<impl Responder, Box<dyn Error>> {
    let mut conn = data.db.get_conn().await?;
    let mut tx = conn.begin().await?;

    let orders = data.db.all_orders(&mut tx).await?;

    tx.commit().await?;
    conn.close().await?;
    Ok(web::Json(orders))
}

#[get("/order/{order_id}")]
async fn get_order(path: web::Path<(Uuid,)>, data: web::Data<AppState>) -> Result<impl Responder, Box<dyn Error>> {
    let mut conn = data.db.get_conn().await?;
    let mut tx = conn.begin().await?;

    let order = data.db.get_order_by_uuid(&mut tx, path.0).await?;

    tx.commit().await?;
    conn.close().await?;
    Ok(web::Json(order))
}

#[get("/order/{order_id}/order-item")]
async fn get_order_items(path: web::Path<(Uuid,)>, data: web::Data<AppState>) -> Result<impl Responder, Box<dyn Error>> {
    let mut conn = data.db.get_conn().await?;
    let mut tx = conn.begin().await?;

    let order_items = data.db.all_order_items(&mut tx, path.0).await?;

    tx.commit().await?;
    conn.close().await?;
    Ok(web::Json(order_items))
}

#[post("/order/{order_id}/order-item")]
async fn create_order_item(new_order_item: web::Json<NewOrderItem>, data: web::Data<AppState>) -> Result<impl Responder, Box<dyn Error>> {
    let mut conn = data.db.get_conn().await?;
    let mut tx = conn.begin().await?;

    let order_item = data.db.create_order_item(&mut tx, new_order_item.into_inner()).await?;

    tx.commit().await?;
    conn.close().await?;
    Ok(web::Json(order_item))
}
