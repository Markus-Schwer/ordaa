use diesel::prelude::*;

use crate::{dto::{NewMenuDto, NewMenuItemDto}, schema::menu_items::menu_id};

#[derive(Identifiable, Queryable, Selectable)]
#[diesel(table_name = crate::schema::menus)]
#[diesel(check_for_backend(diesel::sqlite::Sqlite))]
pub struct Menu {
    pub id: i32,
    pub name: String,
    pub url: Option<String>,
}

#[derive(Insertable)]
#[diesel(table_name = crate::schema::menus)]
#[diesel(check_for_backend(diesel::sqlite::Sqlite))]
pub struct NewMenu {
    pub name: String,
    pub url: Option<String>,
}

impl NewMenu {
    pub fn from_dto(dto: NewMenuDto) -> Self { Self { name: dto.name, url: dto.url } }
}

#[derive(Identifiable, Queryable, Selectable, Associations, Clone)]
#[diesel(table_name = crate::schema::menu_items)]
#[diesel(belongs_to(Menu))]
#[diesel(check_for_backend(diesel::sqlite::Sqlite))]
pub struct MenuItem {
    pub id: i32,
    pub short_name: String,
    pub name: String,
    pub menu_id: i32,
    pub price: i32,
}

#[derive(Insertable)]
#[diesel(table_name = crate::schema::menu_items)]
#[diesel(check_for_backend(diesel::sqlite::Sqlite))]
pub struct NewMenuItem {
    pub short_name: String,
    pub name: String,
    pub menu_id: i32,
    pub price: i32,
}

impl NewMenuItem {
    pub fn from_dto(dto: NewMenuItemDto, menu: i32) -> Self { Self { short_name: dto.short_name, name: dto.name, menu_id: menu, price: dto.price } }
}

#[derive(Identifiable, Queryable, Selectable)]
#[diesel(table_name = crate::schema::users)]
#[diesel(check_for_backend(diesel::sqlite::Sqlite))]
pub struct User {
    pub id: i32,
    pub name: String
}

#[derive(Insertable)]
#[diesel(table_name = crate::schema::users)]
#[diesel(check_for_backend(diesel::sqlite::Sqlite))]
pub struct NewUser {
    pub name: String
}

#[derive(Identifiable, Queryable, Selectable, Associations)]
#[diesel(table_name = crate::schema::orders)]
#[diesel(belongs_to(User, foreign_key = initiator))]
#[diesel(check_for_backend(diesel::sqlite::Sqlite))]
pub struct Order {
    pub id: i32,
    pub order_deadline: Option<i32>,
    pub eta: Option<i32>,
    pub initiator: i32,
    pub sugar_person: Option<i32>,
    pub state: String,
}

#[derive(Insertable)]
#[diesel(table_name = crate::schema::orders)]
#[diesel(check_for_backend(diesel::sqlite::Sqlite))]
pub struct NewOrder {
    pub order_deadline: Option<i32>,
    pub eta: Option<i32>,
    pub initiator: i32,
    pub sugar_person: Option<i32>,
    pub state: String,
}

#[derive(Identifiable, Queryable, Selectable, Associations)]
#[diesel(table_name = crate::schema::order_items)]
#[diesel(belongs_to(Order))]
#[diesel(check_for_backend(diesel::sqlite::Sqlite))]
pub struct OrderItem {
    pub id: i32,
    pub menu_item_id: i32,
    pub user: i32,
    pub paid: bool,
    pub order_id: i32,
}

#[derive(Insertable)]
#[diesel(table_name = crate::schema::order_items)]
#[diesel(check_for_backend(diesel::sqlite::Sqlite))]
pub struct NewOrderItem {
    pub menu_item_id: i32,
    pub user: i32,
    pub paid: bool,
    pub order_id: i32,
}
