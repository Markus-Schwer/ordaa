use serde::{Serialize, Deserialize};

use crate::entity::models::{MenuItem, Menu};

#[derive(Serialize, Deserialize, Clone)]
pub struct MenuDto {
    pub id: i32,
    pub name: String,
    pub url: Option<String>,
    pub items: Vec<MenuItemDto>
}

impl MenuDto {
    pub fn from_db(menu: Menu, items: Vec<MenuItemDto>) -> Self { Self { id: menu.id, name: menu.name, url: menu.url, items } }
}

#[derive(Serialize, Deserialize, Clone)]
pub struct NewMenuDto {
    pub name: String,
    pub url: Option<String>,
    pub items: Vec<NewMenuItemDto>
}

#[derive(Serialize, Deserialize, Clone)]
pub struct MenuItemDto {
    pub id: i32,
    pub short_name: String,
    pub name: String,
    pub price: i32,
}

impl MenuItemDto {
    pub fn from_db(menu_item: MenuItem) -> Self {
        Self { id: menu_item.id, short_name: menu_item.short_name, name: menu_item.name, price: menu_item.price }
    }
}

#[derive(Serialize, Deserialize, Clone)]
pub struct NewMenuItemDto {
    pub short_name: String,
    pub name: String,
    pub price: i32,
}

#[derive(Serialize, Deserialize, Clone)]
pub struct UserDto {
    pub id: i32,
    pub name: String
}

#[derive(Serialize, Deserialize, Clone)]
pub struct NewUserDto {
    pub name: String
}

#[derive(Serialize, Deserialize, Clone)]
pub struct OrderDto {
    pub id: i32,
    pub order_deadline: Option<i32>,
    pub eta: Option<i32>,
    pub initiator: UserDto,
    pub sugar_person: Option<UserDto>,
    pub state: String,
    pub items: Vec<OrderItemDto>
}

#[derive(Serialize, Deserialize, Clone)]
pub struct NewOrderDto {
    pub order_deadline: Option<i32>,
    pub eta: Option<i32>,
    pub initiator: i32,
    pub sugar_person: Option<i32>,
    pub state: String,
}

#[derive(Serialize, Deserialize, Clone)]
pub struct OrderItemDto {
    pub id: i32,
    pub menu_item_id: i32,
    pub user: UserDto,
    pub paid: bool,
    pub order: OrderDto,
}

#[derive(Serialize, Deserialize, Clone)]
pub struct NewOrderItemDto {
    pub menu_item_id: i32,
    pub user: i32,
    pub paid: bool,
    pub order_id: i32,
}
