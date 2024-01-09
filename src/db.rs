use sqlx::{sqlite::SqliteConnectOptions, pool::PoolConnection, Sqlite, SqlitePool, Pool};

use crate::{menu::{MenuItem, Menu}, orders::state::CreateOrderDTO};

#[derive(Clone)]
pub struct Db {
    pool: Pool<Sqlite>
}

impl Db {
    pub async fn new() -> Self {
        let options = SqliteConnectOptions::new()
            .create_if_missing(true)
            .filename("db.sqlite");
            Db { pool: SqlitePool::connect_with(options).await.unwrap() }
    }

    async fn get_conn(&self) -> PoolConnection<Sqlite> {
        return self.pool.acquire().await.unwrap();
    }

    pub async fn init_schema(&self) {
        sqlx::migrate!().run(&self.pool).await.unwrap();
    }

    pub async fn all_items(&self, menu: String) -> Vec<MenuItem> {
        let mut conn = self.get_conn().await;
        sqlx::query_as!(
            MenuItem,
            "SELECT * FROM MENU_ITEM WHERE menu = ?1",
            menu
        )
        .fetch_all(&mut *conn)
        .await
        .unwrap()
    }

    pub async fn get_items_by_id(&self, ids: Vec<String>, menu: String) -> Vec<MenuItem> {
        let mut conn = self.get_conn().await;
        let mut items = Vec::<MenuItem>::new();
        for id in ids {
            items.push(
                sqlx::query_as!(MenuItem,
                                "SELECT * FROM menu_item WHERE id = ?1 AND menu = ?2",
                                id,
                                menu
                               )
                .fetch_one(&mut *conn)
                .await
                .unwrap()
            );
        }
        items
    }

    pub async fn get_item_by_id(&self, id: String, menu: String) -> MenuItem {
        let mut conn = self.get_conn().await;
        sqlx::query_as!(MenuItem,
                "SELECT * FROM menu_item WHERE id = ?1 AND menu = ?2",
                id,
                menu
            )
            .fetch_one(&mut *conn)
            .await
            .unwrap()
    }

    pub async fn insert_menu(&self, menu: Menu) {
        let mut conn = self.get_conn().await;
        for it in menu.items {
            sqlx::query("INSERT INTO MENU_ITEM VALUES (?1, ?2, ?3, ?4)")
                .bind(menu.name.clone().to_lowercase())
                .bind(it.id)
                .bind(it.name)
                .bind(it.price as i64)
                .execute(&mut *conn)
                .await
                .unwrap();
        }
    }

    pub async fn create_order(&self, dto: CreateOrderDTO) {
        let mut conn = self.get_conn().await;
        sqlx::query("INSERT INTO ORDERS VALUES (?1, ?2, ?3, ?4)")
            .bind(menu.name.clone().to_lowercase())
            .bind(it.id)
            .bind(it.name)
            .bind(it.price as i64)
            .execute(&mut *conn)
            .await
            .unwrap();

    }
}
