use sqlx::types::Uuid;
use serde::{Serialize, Deserialize};

#[derive(sqlx::FromRow, sqlx::Type, Serialize, Deserialize, Clone)]
pub struct Menu {
    pub uuid: Uuid,
    pub name: String,
    pub url: Option<String>,
}

#[derive(sqlx::FromRow, sqlx::Type, Clone)]
pub struct NewMenu {
    pub name: String,
    pub url: Option<String>,
}

#[derive(sqlx::FromRow, Serialize, Deserialize, Clone)]
pub struct NewMenuWithItems {
    pub name: String,
    pub url: Option<String>,
    pub items: Vec<NewMenuItem>,
}

#[derive(sqlx::FromRow, sqlx::Type, Serialize, Deserialize, Clone)]
pub struct MenuItem {
    pub uuid: Uuid,
    pub short_name: String,
    pub name: String,
    pub menu_uuid: Uuid,
    pub price: i32,
}

#[derive(sqlx::FromRow, sqlx::Type, Serialize, Deserialize, Clone)]
pub struct NewMenuItem {
    pub short_name: String,
    pub name: String,
    pub price: i32,
}

#[derive(sqlx::FromRow, Serialize, Deserialize, Clone)]
pub struct MenuWithItems {
    pub uuid: Uuid,
    pub name: String,
    pub url: Option<String>,
    pub items: Vec<MenuItem>
}

#[derive(sqlx::FromRow, sqlx::Type, Serialize, Deserialize, PartialEq, Clone)]
pub struct User {
    pub uuid: Uuid,
    pub name: String
}

#[derive(sqlx::FromRow, sqlx::Type, Clone)]
pub struct NewUser {
    pub name: String
}

#[derive(sqlx::FromRow, sqlx::Type, Clone)]
pub struct Order {
    pub uuid: Uuid,
    pub order_deadline: Option<i32>,
    pub eta: Option<i32>,
    pub initiator: i32,
    pub sugar_person: Option<i32>,
    pub state: String,
    pub menu_uuid: Uuid,
}

#[derive(sqlx::FromRow, Serialize, Deserialize, Clone)]
pub struct OrderWithItems {
    pub uuid: Uuid,
    pub order_deadline: Option<i32>,
    pub eta: Option<i32>,
    #[sqlx(json)]
    pub initiator: User,
    #[sqlx(json)]
    pub sugar_person: Option<User>,
    pub state: String,
    #[sqlx(json)]
    pub menu: Menu,
    #[serde(skip)]
    pub items: Vec<OrderItemWithJoins>,
}

#[derive(sqlx::FromRow, sqlx::Type, Serialize, Deserialize, Clone)]
pub struct NewOrder {
    pub order_deadline: Option<i32>,
    pub eta: Option<i32>,
    pub initiator: Uuid,
    pub sugar_person: Option<Uuid>,
    pub state: String,
    pub menu_uuid: Uuid,
}

#[derive(sqlx::FromRow, sqlx::Type, Clone)]
pub struct OrderItem {
    pub uuid: Uuid,
    pub menu_item_uuid: Uuid,
    pub order_user: Uuid,
    pub paid: bool,
    pub order_uuid: Uuid,
    pub price: i32,
}

#[derive(sqlx::FromRow, sqlx::Type, Serialize, Deserialize, Clone)]
pub struct OrderItemWithJoins {
    pub uuid: Uuid,
    #[sqlx(json)]
    pub menu_item: MenuItem,
    #[sqlx(json)]
    pub order_user: User,
    pub paid: bool,
    pub order_uuid: Uuid,
    pub price: i32,
}

#[derive(sqlx::FromRow, sqlx::Type, Serialize, Deserialize, Clone)]
pub struct NewOrderItem {
    pub menu_item_uuid: Uuid,
    pub order_user: Uuid,
    pub order_uuid: Uuid,
    pub price: i32,
}
