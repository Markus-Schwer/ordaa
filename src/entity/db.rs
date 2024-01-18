use diesel::{ExpressionMethods, result};
use diesel_migrations::{embed_migrations, EmbeddedMigrations, MigrationHarness};
pub const MIGRATIONS: EmbeddedMigrations = embed_migrations!("migrations");

use diesel::{sqlite::SqliteConnection, r2d2::Pool};
use diesel::r2d2::{ConnectionManager, PooledConnection};
use diesel::prelude::*;
use diesel::upsert::*;
use std::env;
use itertools::izip;

use crate::boundary::dto::{MenuItemDto, MenuWithItemsDto, NewMenuDto, NewOrderDto, OrderDto, OrderItemDto, UserDto, NewOrderItemDto, MenuDto};
use crate::entity::models::{MenuItem, Menu, Order, NewMenu, NewOrder, NewMenuItem};
use crate::entity::schema::{menus, orders, menu_items};

use super::models::{OrderItem, User, NewOrderItem};
use super::schema::{order_items, users};

#[derive(Clone)]
pub struct Db {
    pool: Pool<ConnectionManager<SqliteConnection>>
}

pub type Connection = PooledConnection<ConnectionManager<SqliteConnection>>;

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

    pub fn get_conn(&self) -> Result<Connection, r2d2::Error> {
        self.pool.get()
    }

    pub fn init_schema(&self) {
        self.get_conn().unwrap().run_pending_migrations(MIGRATIONS).unwrap();
    }

    pub fn all_menus(&self, conn: &mut Connection) -> Result<Vec<MenuWithItemsDto>, result::Error> {
        let menus = menus::table.select(Menu::as_select()).load(conn)?;
        let items = MenuItem::belonging_to(&menus).select(MenuItem::as_select()).load(conn)?;
        Ok(items.grouped_by(&menus).into_iter().zip(menus)
            .map(|(mi, m): (Vec<MenuItem>, Menu)| MenuWithItemsDto::from_db(m, mi.iter().map(|mi: &MenuItem| MenuItemDto::from_db(mi.clone())).collect::<Vec<MenuItemDto>>()))
            .collect::<Vec<MenuWithItemsDto>>())
    }

    pub fn get_menu_by_id(&self, conn: &mut Connection, id: i32) -> Result<MenuWithItemsDto, result::Error> {
        let menu = menus::table.filter(menus::id.eq(id)).select(Menu::as_select()).first(conn)?;
        let items = MenuItem::belonging_to(&menu).select(MenuItem::as_select()).load(conn)?
            .iter().map(|mi: &MenuItem| MenuItemDto::from_db(mi.clone())).collect();
        Ok(MenuWithItemsDto::from_db(menu, items))
    }

    pub fn all_menu_items(&self, conn: &mut Connection, menu_id: i32) -> Result<Vec<MenuItemDto>, result::Error> {
        Ok(menu_items::table.filter(menu_items::menu_id.eq(menu_id)).select(MenuItem::as_select()).load(conn)?
            .iter().map(|m: &MenuItem| MenuItemDto::from_db(m.clone())).collect())
    }

    pub fn get_items_by_id(&self, conn: &mut Connection, ids: Vec<i32>, menu_id: i32) -> Result<Vec<MenuItemDto>, result::Error> {
        Ok(menu_items::table.filter(menu_items::id.eq_any(ids).and(menu_items::menu_id.eq(menu_id))).select(MenuItem::as_select()).load(conn)?
            .iter().map(|m: &MenuItem| MenuItemDto::from_db(m.clone())).collect())
    }

    pub fn get_item_by_id(&self, conn: &mut Connection, id: i32, menu: i32) -> Result<MenuItemDto, result::Error> {
        Ok(MenuItemDto::from_db(menu_items::table.filter(menu_items::menu_id.eq(menu).and(menu_items::id.eq(id))).select(MenuItem::as_select()).first(conn)?))
    }

    pub fn upsert_menu(&self, conn: &mut Connection, id: i32, menu: NewMenuDto) -> Result<MenuWithItemsDto, result::Error> {
        let result = diesel::insert_into(menus::table)
            .values(&Menu::from_dto(id, menu.clone()))
            .on_conflict(menus::id)
            .do_update()
            .set((
                menus::name.eq(menu.name),
                menus::url.eq(menu.url)
            )).returning(Menu::as_returning()).get_result(conn)?;

        let res_id = result.id;
        for menu_item in menu.items {
            let existing_menu_item: Option<MenuItem> = menu_items::table.filter(menu_items::menu_id.eq(res_id).and(menu_items::short_name.eq(menu_item.short_name.to_string())))
                .select(MenuItem::as_select())
                .first(conn).optional()?;
            if let Some(existing) = existing_menu_item {
                diesel::update(menu_items::table.find(existing.id))
                    .set((
                        menu_items::name.eq(menu_item.name),
                        menu_items::price.eq(menu_item.price)
                    )).execute(conn)?;
            } else {
                diesel::insert_into(menu_items::table).values(&NewMenuItem::from_dto(menu_item, res_id)).execute(conn)?;
            }
        }
        Ok(MenuWithItemsDto::from_db(result, self.all_menu_items(conn, res_id)?))
    }

    pub fn insert_menu(&self, conn: &mut Connection, menu: NewMenuDto) -> Result<MenuWithItemsDto, result::Error> {
        let result = diesel::insert_into(menus::table).values(&NewMenu::from_dto(menu.clone())).returning(Menu::as_returning()).get_result(conn)?;
        let res_id = result.id;
        for menu_item in menu.items {
            diesel::insert_into(menu_items::table).values(&NewMenuItem::from_dto(menu_item, res_id)).execute(conn)?;
        }
        Ok(MenuWithItemsDto::from_db(result, self.all_menu_items(conn, res_id)?))
    }

    pub fn all_orders(&self, conn: &mut Connection) -> Result<Vec<OrderDto>, result::Error> {
        let orders = orders::table.select(Order::as_select()).load(conn)?;
        let items = OrderItem::belonging_to(&orders)
            .inner_join(users::table)
            .inner_join(menu_items::table)
            .select((OrderItem::as_select(), User::as_select(), MenuItem::as_select()))
            .load::<(OrderItem, User, MenuItem)>(conn)?
            .grouped_by(&orders);

        let initiators = orders::table
            .inner_join(users::table)
            .select(User::as_select())
            .load(conn)?;

        let sugar_persons = orders::table
            .left_join(users::table)
            .select(Option::<User>::as_select())
            .load(conn)?;

        let joined_menus = orders::table
            .inner_join(menus::table)
            .select(Menu::as_select())
            .load(conn)?;

        Ok(izip!(&items, &initiators, &sugar_persons, &joined_menus, &orders).map(|(oiu, i, sp, m, o): (&Vec<(OrderItem, User, MenuItem)>, &User, &Option<User>, &Menu, &Order)| {
            OrderDto::from_db(
                o.clone(),
                UserDto::from_db(i.clone()),
                sp.clone().map(|u| UserDto::from_db(u)),
                MenuDto::from_db(m.clone()),
                oiu.iter().map(|(oi, u, mi): &(OrderItem, User, MenuItem)| OrderItemDto::from_db(oi.clone(), UserDto::from_db(u.clone()), MenuItemDto::from_db(mi.clone()))).collect::<Vec<OrderItemDto>>()
                )
        }).collect::<Vec<OrderDto>>())
    }

    pub fn get_order_by_id(&self, conn: &mut Connection, id: i32) -> Result<OrderDto, result::Error> {
        let order = orders::table.filter(orders::id.eq(id)).select(Order::as_select()).first(conn)?;
        let order_id = order.id;

        let initiator = UserDto::from_db(orders::table.find(order_id)
            .inner_join(users::table)
            .select(User::as_select())
            .first(conn)?);

        // TODO: diesel currently only supports joining by target table, not by foreign key, which
        // means that two foreign keys on orders to users cannot be queried. For now, only
        // initiator can be queried. The following query should fetch 'sugar_person', but only
        // fetches initiator.
        let sugar_person: Option<UserDto> = orders::table.find(order_id)
            .left_join(users::table)
            .select(Option::<User>::as_select())
            .first::<Option<User>>(conn)?
            .map(|u| UserDto::from_db(u));

        let menu = MenuDto::from_db(orders::table.find(order_id)
                                    .inner_join(menus::table)
                                    .select(Menu::as_select())
                                    .first(conn)?);

        let items = self.all_order_items(conn, order.id)?;

        Ok(OrderDto::from_db(order, initiator, sugar_person, menu, items))
    }

    pub fn all_order_items(&self, conn: &mut Connection, order_id: i32) -> Result<Vec<OrderItemDto>, result::Error> {
        Ok(order_items::table
            .inner_join(users::table)
            .inner_join(menu_items::table)
            .filter(order_items::order_id.eq(order_id))
            .select((OrderItem::as_select(), User::as_select(), MenuItem::as_select()))
            .load::<(OrderItem, User, MenuItem)>(conn)?
            .iter().map(|(oi, u, mi): &(OrderItem, User, MenuItem)| OrderItemDto::from_db(oi.clone(), UserDto::from_db(u.clone()), MenuItemDto::from_db(mi.clone()))).collect())
    }

    pub fn create_order(&self, conn: &mut Connection, order: NewOrderDto) -> Result<OrderDto, result::Error> {
        let order = diesel::insert_into(orders::table)
            .values(&NewOrder::from_dto(order))
            .returning(Order::as_returning())
            .get_result(conn)?;
        let order_id = order.id;

        let initiator = UserDto::from_db(orders::table.find(order_id)
            .inner_join(users::table)
            .select(User::as_select())
            .first(conn)?);

        // TODO: diesel currently only supports joining by target table, not by foreign key, which
        // means that two foreign keys on orders to users cannot be queried. For now, only
        // initiator can be queried. The following query should fetch 'sugar_person', but only
        // fetches initiator.
        let sugar_person: Option<UserDto> = orders::table.find(order_id)
            .left_join(users::table)
            .select(Option::<User>::as_select())
            .first::<Option<User>>(conn)?
            .map(|u| UserDto::from_db(u));

        let menu = MenuDto::from_db(orders::table.find(order_id)
                                    .inner_join(menus::table)
                                    .select(Menu::as_select())
                                    .first(conn)?);

        Ok(OrderDto::from_db(order, initiator, sugar_person, menu, self.all_order_items(conn, order_id)?))
    }

    pub fn create_order_item(&self, conn: &mut Connection, order_item: NewOrderItemDto) -> Result<OrderItemDto, result::Error> {
        let price: i32 = menu_items::table.find(order_item.menu_item_id).select(menu_items::price).first(conn)?;
        let order_item = diesel::insert_into(order_items::table)
            .values(&NewOrderItem::from_dto(order_item, price))
            .returning(OrderItem::as_returning())
            .get_result(conn)?;
        let order_item_id = order_item.id;

        let user = order_items::table.find(order_item_id)
            .inner_join(users::table)
            .select(User::as_select())
            .first(conn)?;

        let menu_item = order_items::table.find(order_item_id)
            .inner_join(menu_items::table)
            .select(MenuItem::as_select())
            .first(conn)?;

        Ok(OrderItemDto::from_db(order_item, UserDto::from_db(user), MenuItemDto::from_db(menu_item)))
    }
}
