use sqlx::{sqlite::SqliteConnectOptions, pool::PoolConnection, Sqlite, SqlitePool, Pool};

use crate::menu::{MenuItem, Menu};

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

    pub async fn get_items_by_id(&self, maybe_ids: Option<Vec<String>>, menu: String) -> Vec<MenuItem> {
        let mut conn = self.get_conn().await;
        let mut matches: Vec<MenuItem> = Vec::new();
        if let Some(ids) = maybe_ids {
            for id in ids {
                matches.push(
                    sqlx::query_as::<_, MenuItem>(
                        "SELECT id, name, price FROM MENU_ITEM WHERE id = ?1 AND menu = ?2",
                    )
                    .bind(id)
                    .bind(menu.clone())
                    .fetch_one(&mut *conn)
                    .await
                    .unwrap(),
                );
            }
        } else {
            matches = sqlx::query_as::<_, MenuItem>(
                "SELECT id, name, price FROM MENU_ITEM WHERE menu = ?1",
            )
            .bind(menu)
            .fetch_all(&mut *conn)
            .await
            .unwrap()
        }
        matches
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
}
