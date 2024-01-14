use std::error::Error;

use actix_web::{get, put, web, HttpResponse, Responder};
use serde::{Deserialize, Serialize};
use crate::{boundary::dto::MenuItemDto, service::state::AppState};

use super::dto::NewMenuDto;

pub struct Menu {
    pub id: i32,
    pub name: String,
    pub url: Option<String>,
    pub items: Vec<MenuItemDto>
}

pub struct MenuItem {
    pub id: i32,
    pub short_name: String,
    pub name: String,
    pub price: i32,
}

pub fn services_menu(cfg: &mut web::ServiceConfig) {
    cfg.service(get_menu);
    cfg.service(put_menu);
}

#[put("/menu/{menu_id}")]
async fn put_menu(_: web::Path<(i32,)>, new_menu: web::Json<NewMenuDto>, data: web::Data<AppState>) -> Result<HttpResponse, Box<dyn Error>> {
    let menu = data.db.insert_menu(new_menu.into_inner());
    data.search.index_write_sender.send(menu).unwrap();
    Ok(HttpResponse::Ok().finish())
}

#[derive(Deserialize, Serialize, Debug)]
struct FuzzyParam {
    search_string: String,
}

#[get("/menu/{menu_id}")]
async fn get_menu(path: web::Path<(i32,)>, query: web::Query<Option<FuzzyParam>>, data: web::Data<AppState>) -> Result<impl Responder, Box<dyn Error>> {
    let items = if let Some(param) = query.into_inner() {
        let ids = data.search.fuzz_menu_item_ids(param.search_string.as_str());
        data.db.get_items_by_id(ids, path.0)
    } else {
        data.db.all_items(path.0)
    };
    Ok(web::Json(items))
}
