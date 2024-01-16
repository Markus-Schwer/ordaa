use diesel::ExpressionMethods;
use diesel_migrations::{embed_migrations, EmbeddedMigrations, MigrationHarness};
pub const MIGRATIONS: EmbeddedMigrations = embed_migrations!("migrations");

use diesel::{sqlite::SqliteConnection, r2d2::Pool};
use diesel::r2d2::{ConnectionManager, PooledConnection};
use diesel::prelude::*;
use diesel::upsert::*;
use std::env;
use itertools::izip;

use crate::boundary::dto::{MenuItemDto, MenuDto, NewMenuDto, NewMenuItemDto, NewOrderDto, OrderDto, OrderItemDto, UserDto};
use crate::entity::models::{MenuItem, Menu, Order, NewMenu, NewOrder, NewMenuItem};
use crate::entity::schema::{menus, orders, menu_items};

use super::models::{OrderItem, User};
use super::schema::{order_items, users};

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

    pub fn all_menu_items(&self, menu_id: i32) -> Vec<MenuItemDto> {
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

    pub fn upsert_menu(&self, id: i32, menu: NewMenuDto) -> MenuDto {
        let mut conn = self.get_conn();

        let result = diesel::insert_into(menus::table)
            .values(&Menu::from_dto(id, menu.clone()))
            .on_conflict(menus::id)
            .do_update()
            .set((
                menus::name.eq(menu.name),
                menus::url.eq(menu.url)
            )).returning(Menu::as_returning()).get_result(&mut conn).unwrap();

        let res_id = result.id;
        for menu_item in menu.items {
            let existing_menu_item: Option<MenuItem> = menu_items::table.filter(menu_items::menu_id.eq(res_id).and(menu_items::short_name.eq(menu_item.short_name.to_string())))
                .select(MenuItem::as_select())
                .first(&mut conn).optional().unwrap();
            if let Some(existing) = existing_menu_item {
                diesel::update(menu_items::table.find(existing.id))
                    .set((
                        menu_items::name.eq(menu_item.name),
                        menu_items::price.eq(menu_item.price)
                    )).execute(&mut conn).unwrap();
            } else {
                diesel::insert_into(menu_items::table).values(&NewMenuItem::from_dto(menu_item, res_id)).execute(&mut conn).unwrap();
            }
        }
        MenuDto::from_db(result, self.all_menu_items(res_id))
    }

    pub fn insert_menu(&self, menu: NewMenuDto) -> MenuDto {
        let mut conn = self.get_conn();
        let result = diesel::insert_into(menus::table).values(&NewMenu::from_dto(menu.clone())).returning(Menu::as_returning()).get_result(&mut conn).unwrap();
        let res_id = result.id;
        for menu_item in menu.items {
            diesel::insert_into(menu_items::table).values(&NewMenuItem::from_dto(menu_item, res_id)).execute(&mut conn).unwrap();
        }
        MenuDto::from_db(result, self.all_menu_items(res_id))
    }

    pub fn all_orders(&self) -> Vec<OrderDto> {
        let mut conn = self.get_conn();
        let orders = orders::table.select(Order::as_select()).load(&mut conn).unwrap();
        let items = OrderItem::belonging_to(&orders)
            .inner_join(users::table)
            .select((OrderItem::as_select(), User::as_select()))
            .load::<(OrderItem, User)>(&mut conn).unwrap()
            .grouped_by(&orders);

        let initiators = orders::table
            .inner_join(users::table)
            .select(User::as_select())
            .load(&mut conn)
            .unwrap();

        let sugar_persons = orders::table
            .left_join(users::table)
            .select(Option::<User>::as_select())
            .load(&mut conn)
            .unwrap();

        izip!(&items, &initiators, &sugar_persons, &orders).map(|(oiu, i, sp, o): (&Vec<(OrderItem, User)>, &User, &Option<User>, &Order)| {
            OrderDto::from_db(
                o.clone(),
                UserDto::from_db(i.clone()),
                sp.clone().map(|u| UserDto::from_db(u)),
                oiu.iter().map(|(oi, u): &(OrderItem, User)| OrderItemDto::from_db(oi.clone(), UserDto::from_db(u.clone()))).collect::<Vec<OrderItemDto>>()
                )
        }).collect::<Vec<OrderDto>>()
    }

    pub fn all_order_items(&self, order_id: i32) -> Vec<OrderItemDto> {
        let mut conn = self.get_conn();

        order_items::table
            .inner_join(users::table)
            .filter(order_items::order_id.eq(order_id))
            .select((OrderItem::as_select(), User::as_select()))
            .load::<(OrderItem, User)>(&mut conn).unwrap()
            .iter().map(|(oi, u): &(OrderItem, User)| OrderItemDto::from_db(oi.clone(), UserDto::from_db(u.clone()))).collect()
    }

    pub fn create_order(&self, order: NewOrderDto) -> OrderDto {
        let mut conn = self.get_conn();
        let order = diesel::insert_into(orders::table)
            .values(&NewOrder::from_dto(order))
            .returning(Order::as_returning())
            .get_result(&mut conn).unwrap();
        let order_id = order.id;

        let initiator = orders::table.find(order_id)
            .inner_join(users::table)
            .select(User::as_select())
            .first(&mut conn)
            .unwrap();

        // TODO: diesel currently only supports joining by target table, not by foreign key, which
        // means that two foreign keys on orders to users cannot be queried. For now, only
        // initiator can be queried. The following query should fetch 'sugar_person', but only
        // fetches initiator.
        let sugar_person: Option<User> = orders::table.find(order_id)
            .left_join(users::table)
            .select(Option::<User>::as_select())
            .first::<Option<User>>(&mut conn)
            .unwrap();
        OrderDto::from_db(order, UserDto::from_db(initiator), sugar_person.map(|u| UserDto::from_db(u)), self.all_order_items(order_id))
    }
}
