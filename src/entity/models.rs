use diesel::prelude::*;

use crate::boundary::dto::{NewMenuDto, NewMenuItemDto, NewOrderDto, NewOrderItemDto};

#[derive(Identifiable, Queryable, Selectable, Insertable)]
#[diesel(table_name = crate::entity::schema::menus)]
#[diesel(check_for_backend(diesel::sqlite::Sqlite))]
pub struct Menu {
    pub id: i32,
    pub name: String,
    pub url: Option<String>,
}

impl Menu {
    pub fn from_dto(id: i32, dto: NewMenuDto) -> Self { Self { id, name: dto.name, url: dto.url } }
}

#[derive(Insertable)]
#[diesel(table_name = crate::entity::schema::menus)]
#[diesel(check_for_backend(diesel::sqlite::Sqlite))]
pub struct NewMenu {
    pub name: String,
    pub url: Option<String>,
}

impl NewMenu {
    pub fn from_dto(dto: NewMenuDto) -> Self { Self { name: dto.name, url: dto.url } }
}

#[derive(Identifiable, Queryable, Selectable, Associations, Clone)]
#[diesel(table_name = crate::entity::schema::menu_items)]
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
#[diesel(table_name = crate::entity::schema::menu_items)]
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

#[derive(Identifiable, Queryable, Selectable, Clone)]
#[diesel(table_name = crate::entity::schema::users)]
#[diesel(check_for_backend(diesel::sqlite::Sqlite))]
pub struct User {
    pub id: i32,
    pub name: String
}

#[derive(Insertable)]
#[diesel(table_name = crate::entity::schema::users)]
#[diesel(check_for_backend(diesel::sqlite::Sqlite))]
pub struct NewUser {
    pub name: String
}

#[derive(Identifiable, Queryable, Selectable, Associations, Clone)]
#[diesel(table_name = crate::entity::schema::orders)]
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
#[diesel(table_name = crate::entity::schema::orders)]
#[diesel(check_for_backend(diesel::sqlite::Sqlite))]
pub struct NewOrder {
    pub order_deadline: Option<i32>,
    pub eta: Option<i32>,
    pub initiator: i32,
    pub sugar_person: Option<i32>,
    pub state: String,
}

impl NewOrder {
    pub fn from_dto(dto: NewOrderDto) -> Self { Self { order_deadline: dto.order_deadline, eta: dto.eta, initiator: dto.initiator, sugar_person: dto.sugar_person, state: dto.state } }
}

#[derive(Identifiable, Queryable, Selectable, Associations, Clone)]
#[diesel(table_name = crate::entity::schema::order_items)]
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
#[diesel(table_name = crate::entity::schema::order_items)]
#[diesel(check_for_backend(diesel::sqlite::Sqlite))]
pub struct NewOrderItem {
    pub menu_item_id: i32,
    pub user: i32,
    pub order_id: i32,
}

impl NewOrderItem {
    pub fn from_dto(dto: NewOrderItemDto) -> Self { Self { menu_item_id: dto.menu_item_id, user: dto.user, order_id: dto.order_id } }
}
