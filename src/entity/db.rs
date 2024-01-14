use diesel::ExpressionMethods;
use diesel_migrations::{embed_migrations, EmbeddedMigrations, MigrationHarness};
pub const MIGRATIONS: EmbeddedMigrations = embed_migrations!("migrations");

use diesel::{sqlite::SqliteConnection, r2d2::Pool};
use diesel::r2d2::{ConnectionManager, PooledConnection};
use diesel::prelude::*;
use std::env;

use crate::boundary::dto::{MenuItemDto, MenuDto, NewMenuDto};
use crate::entity::models::{MenuItem, Menu, Order, NewMenu, NewOrder, NewMenuItem};
use crate::entity::schema::{menus, orders, menu_items};

#[derive(Clone)]
pub struct Db {
    pool: Pool<ConnectionManager<SqliteConnection>>
}

impl Db {
    pub fn new() -> Self {
        let database_url = env::var("DATABASE_URL").expect("DATABASE_URL must be set");
        let manager = ConnectionManager::<SqliteConnection>::new(database_url);
        // Refer to the `r2d2` documentation for more methods to use
        // when building a connection pool
        let pool = Pool::builder()
            .test_on_check_out(true)
            .build(manager)
            .expect("Could not build connection pool");
        Db { pool }
    }

    fn get_conn(&self) -> PooledConnection<ConnectionManager<SqliteConnection>> {
        return self.pool.get().unwrap();
    }

    pub fn init_schema(&self) {
        self.get_conn().run_pending_migrations(MIGRATIONS).unwrap();
    }

    pub fn all_menus(&self) -> Vec<MenuDto> {
        let mut conn = self.get_conn();
        let menus = menus::table.select(Menu::as_select()).load(&mut conn).unwrap();
        let items = MenuItem::belonging_to(&menus).select(MenuItem::as_select()).load(&mut conn).unwrap();
        items.grouped_by(&menus).into_iter().zip(menus)
            .map(|(mi, m): (Vec<MenuItem>, Menu)| MenuDto::from_db(m, mi.iter().map(|mi: &MenuItem| MenuItemDto::from_db(mi.clone())).collect::<Vec<MenuItemDto>>()))
            .collect::<Vec<MenuDto>>()
    }

    pub fn get_menu_by_id(&self, id: i32) -> MenuDto {
        let mut conn = self.get_conn();

        let menu = menus::table.filter(menus::id.eq(id)).select(Menu::as_select()).first(&mut conn).unwrap();
        let items = MenuItem::belonging_to(&menu).select(MenuItem::as_select()).load(&mut conn).unwrap()
            .iter().map(|mi: &MenuItem| MenuItemDto::from_db(mi.clone())).collect();
        MenuDto::from_db(menu, items)
    }

    pub fn all_items(&self, menu_id: i32) -> Vec<MenuItemDto> {
        let mut conn = self.get_conn();
        menu_items::table.filter(menu_items::menu_id.eq(menu_id)).select(MenuItem::as_select()).load(&mut conn).unwrap()
            .iter().map(|m: &MenuItem| MenuItemDto::from_db(m.clone())).collect()
    }

    pub fn get_items_by_id(&self, ids: Vec<i32>, menu_id: i32) -> Vec<MenuItemDto> {
        let mut conn = self.get_conn();

        menu_items::table.filter(menu_items::id.eq_any(ids).and(menu_items::menu_id.eq(menu_id))).select(MenuItem::as_select()).load(&mut conn).unwrap()
            .iter().map(|m: &MenuItem| MenuItemDto::from_db(m.clone())).collect()
    }

    pub fn get_item_by_id(&self, id: i32, menu: i32) -> MenuItemDto {
        let mut conn = self.get_conn();

        MenuItemDto::from_db(menu_items::table.filter(menu_items::menu_id.eq(menu).and(menu_items::id.eq(id))).select(MenuItem::as_select()).first(&mut conn).unwrap())
    }

    pub fn insert_menu(&self, menu: NewMenuDto) -> MenuDto {
        let mut conn = self.get_conn();
        let result = diesel::insert_into(menus::table).values(&NewMenu::from_dto(menu.clone())).returning(Menu::as_returning()).get_result(&mut conn).unwrap();
        let res_id = result.id;
        for menu_item in menu.items {
            diesel::insert_into(menu_items::table).values(&NewMenuItem::from_dto(menu_item, res_id)).execute(&mut conn).unwrap();
        }
        MenuDto::from_db(result, self.all_items(res_id))
    }

    pub fn create_order(&self, order: NewOrder) {
        let mut conn = self.get_conn();
        diesel::insert_into(orders::table).values(&order).returning(Order::as_returning()).get_result(&mut conn).unwrap();
    }
}
