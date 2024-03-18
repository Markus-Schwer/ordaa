use std::env;
use sqlx::Transaction as SqlxTransaction;
use sqlx::pool::PoolConnection;
use sqlx::{Pool, postgres::Postgres, types::Uuid};

use crate::entity::models::{MenuItem, Menu, NewOrder, MenuWithItems};

use super::models::{User, NewOrderItem, OrderWithItems, OrderItemWithJoins, NewMenuWithItems};

#[derive(Clone)]
pub struct Db {
    pool: Pool<Postgres>
}

pub type Connection = PoolConnection<Postgres>;
pub type Transaction<'a> = SqlxTransaction<'a, Postgres>;

impl Db {
    pub async fn new() -> Self {
        let database_url = env::var("DATABASE_URL").expect("DATABASE_URL must be set");
        let pool = Pool::<Postgres>::connect(&database_url).await.unwrap();
        Db { pool }
    }

    pub async fn get_conn(&self) -> Result<Connection, sqlx::Error> {
        self.pool.acquire().await
    }

    pub async fn init_schema(&self) -> Result<(), sqlx::Error> {
        sqlx::migrate!().run(&self.pool).await?;
        Ok(())
    }

    pub async fn all_menus(&self, tx: &mut Transaction<'_>) -> Result<Vec<MenuWithItems>, sqlx::Error> {
        sqlx::query_as!(
            MenuWithItems,
            r#"
                select m.uuid, m.name, m.url, ARRAY_AGG((mi.uuid, mi.short_name, mi.name, mi.price)) as "items!: Vec<MenuItem>"
                from menus m join menu_items mi on m.uuid = mi.menu_uuid
                group by m.uuid, m.name, m.url;
            "#,
        )
        .fetch_all(&mut **tx)
        .await
    }

    pub async fn get_menu_by_uuid(&self, tx: &mut Transaction<'_>, uuid: Uuid) -> Result<MenuWithItems, sqlx::Error> {
        sqlx::query_as!(
            MenuWithItems,
            r#"
                select m.uuid, m.name, m.url, ARRAY_AGG((mi.uuid, mi.short_name, mi.name, mi.price)) as "items!: Vec<MenuItem>"
                from menus m join menu_items mi on m.uuid = mi.menu_uuid
                where m.uuid = $1
                group by m.uuid, m.name, m.url;
            "#,
            uuid
        )
        .fetch_one(&mut **tx)
        .await
    }

    pub async fn all_menu_items(&self, tx: &mut Transaction<'_>, menu_uuid: Uuid) -> Result<Vec<MenuItem>, sqlx::Error> {
        sqlx::query_as!(
            MenuItem,
            r#"
                select uuid, short_name, name, price, menu_uuid
                from menu_items where menu_uuid = $1;
            "#,
            menu_uuid
        )
        .fetch_all(&mut **tx)
        .await
    }

    pub async fn get_menu_items_by_uuid(&self, tx: &mut Transaction<'_>, uuids: Vec<Uuid>, menu_uuid: Uuid) -> Result<Vec<MenuItem>, sqlx::Error> {
        sqlx::query_as!(
            MenuItem,
            r#"
                select uuid, short_name, name, price, menu_uuid
                from menu_items where menu_uuid = $1 and uuid = any($2);
            "#,
            menu_uuid,
            &uuids
        )
        .fetch_all(&mut **tx)
        .await
    }

    pub async fn get_menu_item_by_uuid(&self, tx: &mut Transaction<'_>, uuid: Uuid, menu_uuid: Uuid) -> Result<MenuItem, sqlx::Error> {
        sqlx::query_as!(
            MenuItem,
            r#"
                select uuid, short_name, name, price, menu_uuid
                from menu_items where uuid = $1 and menu_uuid = $2;
            "#,
            uuid,
            menu_uuid
        )
        .fetch_one(&mut **tx)
        .await
    }

    pub async fn update_menu(&self, tx: &mut Transaction<'_>, uuid: Uuid, menu: NewMenuWithItems) -> Result<MenuWithItems, sqlx::Error> {
        sqlx::query!(
            r#"
                update menus
                set name = $2, url = $3
                where uuid = $1;
            "#,
            uuid,
            menu.name,
            menu.url
        ).execute(&mut **tx).await?;

        for menu_item in menu.items {
            let existing_item = sqlx::query_as!(
                MenuItem,
                r#"
                    select mi.* from menu_items mi where mi.menu_uuid = $1 and mi.short_name = $2;
                "#,
                uuid,
                menu_item.short_name
            ).fetch_optional(&mut **tx).await?;
            if let Some(existing_item) = existing_item {
                sqlx::query!(
                    r#"
                        update menu_items
                        set short_name = $2, name = $3, price = $4
                        where uuid = $1;
                    "#,
                    existing_item.uuid,
                    menu_item.short_name,
                    menu_item.name,
                    menu_item.price
                ).execute(&mut **tx).await?;
            } else {
                sqlx::query!(
                    r#"
                        insert into menu_items (menu_uuid, short_name, name, price)
                        values ($1, $2, $3, $4);
                    "#,
                    uuid,
                    menu_item.short_name,
                    menu_item.name,
                    menu_item.price
                ).execute(&mut **tx).await?;
            }
        }
        self.get_menu_by_uuid(tx, uuid).await
    }

    pub async fn insert_menu(&self, tx: &mut Transaction<'_>, menu: NewMenuWithItems) -> Result<MenuWithItems, sqlx::Error> {
        let menu_uuid = sqlx::query!(
            r#"
                insert into menus (name, url)
                values ($1, $2)
                returning uuid;
            "#,
            menu.name,
            menu.url
        ).fetch_one(&mut **tx).await?.uuid;
        for menu_item in menu.items {
            sqlx::query!(
                r#"
                    insert into menu_items (menu_uuid, short_name, name, price)
                    values ($1, $2, $3, $4);
                "#,
                menu_uuid,
                menu_item.short_name,
                menu_item.name,
                menu_item.price
            ).execute(&mut **tx).await?;
        }
        self.get_menu_by_uuid(tx, menu_uuid).await
    }

    pub async fn all_orders(&self, tx: &mut Transaction<'_>) -> Result<Vec<OrderWithItems>, sqlx::Error> {
        sqlx::query_as!(
            OrderWithItems,
            r#"
                select o.uuid, o.order_deadline, o.eta, o.state,
                    to_json(i) as "initiator!: User",
                    to_json(sp) as "sugar_person: User",
                    to_json(m) as "menu!: Menu",
                    ARRAY_AGG((oi.uuid, oi.paid, oi.price, to_json(oiu), to_json(oim))) as "items!: Vec<OrderItemWithJoins>"
                from orders o
                    join order_items oi on o.uuid = oi.order_uuid
                    join users i on i.uuid = o.initiator
                    join users sp on sp.uuid = o.sugar_person
                    join menus m on m.uuid = o.menu_uuid
                    join users oiu on oiu.uuid = oi.order_user
                    join menu_items oim on oim.uuid = oi.menu_item_uuid
                group by o.uuid, o.order_deadline, o.eta, o.state, o.menu_uuid, i.*, sp.*, m.*;
            "#,
        )
        .fetch_all(&mut **tx)
        .await
    }

    pub async fn get_order_by_uuid(&self, tx: &mut Transaction<'_>, uuid: Uuid) -> Result<OrderWithItems, sqlx::Error> {
        sqlx::query_as!(
            OrderWithItems,
            r#"
                select o.uuid, o.order_deadline, o.eta, o.state,
                    to_json(i) as "initiator!: User",
                    to_json(sp) as "sugar_person: User",
                    to_json(m) as "menu!: Menu",
                    ARRAY_AGG((oi.uuid, oi.paid, oi.price, to_json(oiu), to_json(oim))) as "items!: Vec<OrderItemWithJoins>"
                from orders o
                    join order_items oi on o.uuid = oi.order_uuid
                    join users i on i.uuid = o.initiator
                    join users sp on sp.uuid = o.sugar_person
                    join menus m on m.uuid = o.menu_uuid
                    join users oiu on oiu.uuid = oi.order_user
                    join menu_items oim on oim.uuid = oi.menu_item_uuid
                where o.uuid = $1
                group by o.uuid, o.order_deadline, o.eta, o.state, o.menu_uuid, i.*, sp.*, m.*;
            "#,
            uuid
        )
        .fetch_one(&mut **tx)
        .await
    }

    pub async fn all_order_items(&self, tx: &mut Transaction<'_>, order_uuid: Uuid) -> Result<Vec<OrderItemWithJoins>, sqlx::Error> {
        sqlx::query_as!(
            OrderItemWithJoins,
            r#"
                select oi.uuid, oi.paid, oi.price, oi.order_uuid,
                    to_json(mi) as "menu_item!: MenuItem",
                    to_json(u) as "order_user!: User"
                from order_items oi join users u on u.uuid = oi.order_user join menu_items mi on mi.uuid = oi.menu_item_uuid
                where oi.order_uuid = $1;
            "#,
            order_uuid
        )
        .fetch_all(&mut **tx)
        .await
    }

    pub async fn get_order_item_by_uuid(&self, tx: &mut Transaction<'_>, uuid: Uuid) -> Result<OrderItemWithJoins, sqlx::Error> {
        sqlx::query_as!(
            OrderItemWithJoins,
            r#"
                select oi.uuid, oi.paid, oi.price, oi.order_uuid,
                    to_json(mi) as "menu_item!: MenuItem",
                    to_json(u) as "order_user!: User"
                from order_items oi join users u on u.uuid = oi.order_user join menu_items mi on mi.uuid = oi.menu_item_uuid
                where oi.uuid = $1;
            "#,
            uuid
        )
        .fetch_one(&mut **tx)
        .await
    }

    pub async fn create_order(&self, tx: &mut Transaction<'_>, order: NewOrder) -> Result<OrderWithItems, sqlx::Error> {
        let order_uuid = sqlx::query!(
            r#"
                insert into orders (order_deadline, eta, initiator, sugar_person, state, menu_uuid)
                values ($1, $2, $3, $4, $5, $6)
                returning uuid;
            "#,
            order.order_deadline,
            order.eta,
            order.initiator,
            order.sugar_person,
            order.state,
            order.menu_uuid
        ).fetch_one(&mut **tx).await?.uuid;
        self.get_order_by_uuid(tx, order_uuid).await
    }

    pub async fn create_order_item(&self, tx: &mut Transaction<'_>, order_item: NewOrderItem) -> Result<OrderItemWithJoins, sqlx::Error> {
        let order_item_uuid = sqlx::query!(
            r#"
                insert into order_items (menu_item_uuid, order_user, order_uuid, price)
                values ($1, $2, $3, $4)
                returning uuid;
            "#,
            order_item.menu_item_uuid,
            order_item.order_user,
            order_item.order_uuid,
            order_item.price
        ).fetch_one(&mut **tx).await?.uuid;
        self.get_order_item_by_uuid(tx, order_item_uuid).await
    }

    pub async fn get_user(&self, tx: &mut Transaction<'_>, uuid: Uuid) -> Result<User, sqlx::Error> {
        sqlx::query_as!(
            User,
            r#"
                select uuid, name
                from users where uuid = $1;
            "#,
            uuid
        )
        .fetch_one(&mut **tx)
        .await
    }
}
