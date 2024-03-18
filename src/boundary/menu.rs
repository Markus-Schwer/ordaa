use std::error::Error;

use uuid::Uuid;
use actix_web::{get, post, put, web, Responder};
use serde::{Deserialize, Serialize};
use sqlx::Acquire;
use crate::{service::state::AppState, entity::models::{NewMenuWithItems, NewMenu}};

pub struct Menu {
    pub id: i32,
    pub name: String,
    pub url: Option<String>,
    pub items: Vec<MenuItem>
}

pub struct MenuItem {
    pub id: i32,
    pub short_name: String,
    pub name: String,
    pub price: i32,
}

pub fn services_menu(cfg: &mut web::ServiceConfig) {
    cfg.service(get_menu);
    cfg.service(all_menus);
    cfg.service(create_menu);
    cfg.service(put_menu);
}

#[post("/menu")]
async fn create_menu(new_menu: web::Json<NewMenuWithItems>, data: web::Data<AppState>) -> Result<impl Responder, Box<dyn Error>> {
    let mut conn = data.db.get_conn().await?;
    let mut tx = conn.begin().await?;

    let menu = data.db.insert_menu(&mut tx, new_menu.into_inner()).await?;
    data.search.index_write_sender.send(menu.clone())?;

    tx.commit().await?;
    conn.close().await?;
    Ok(web::Json(menu))
}

#[put("/menu/{menu_uuid}")]
async fn put_menu(path: web::Path<(Uuid,)>, new_menu: web::Json<NewMenuWithItems>, data: web::Data<AppState>) -> Result<impl Responder, Box<dyn Error>> {
    let mut conn = data.db.get_conn().await?;
    let mut tx = conn.begin().await?;

    let menu = data.db.update_menu(&mut tx, path.0, new_menu.into_inner()).await?;

    tx.commit().await?;
    conn.close().await?;

    // TODO: update search index
    // data.search.index_write_sender.send(menu)?;
    Ok(web::Json(menu))
}

#[derive(Deserialize, Serialize, Debug)]
struct FuzzyParam {
    search_string: String,
}

#[get("/menu/{menu_id}")]
async fn get_menu(path: web::Path<(Uuid,)>, query: web::Query<Option<FuzzyParam>>, data: web::Data<AppState>) -> Result<impl Responder, Box<dyn Error>> {
    let mut conn = data.db.get_conn().await?;
    let mut tx = conn.begin().await?;

    let items = if let Some(param) = query.into_inner() {
        let ids = data.search.fuzz_menu_item_ids(param.search_string.as_str());
        data.db.get_menu_items_by_uuid(&mut tx, ids, path.0).await?
    } else {
        data.db.all_menu_items(&mut tx, path.0).await?
    };

    tx.rollback().await?;
    conn.close().await?;
    Ok(web::Json(items))
}

#[get("/menu")]
async fn all_menus(data: web::Data<AppState>) -> Result<impl Responder, Box<dyn Error>> {
    let mut conn = data.db.get_conn().await?;
    let mut tx = conn.begin().await?;

    let menus = web::Json(data.db.all_menus(&mut tx).await?);

    tx.rollback().await?;
    conn.close().await?;
    Ok(menus)
}
