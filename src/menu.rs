use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize, sqlx::FromRow)]
pub struct MenuItem {
    id: String,
    name: String,
    price: i64,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct Menu {
    pub items: Vec<MenuItem>,
}

pub mod filters {
    use sqlx::SqlitePool;
    use warp::{http::StatusCode, Filter};

    use crate::filters::{json_body, with_db};

    use super::{Menu, MenuItem};

    pub fn all(
        pool: &SqlitePool,
    ) -> impl Filter<Extract = (impl warp::Reply,), Error = warp::Rejection> + Clone {
        update_menu(pool).or(get_menu(pool))
    }

    fn get_menu(
        pool: &SqlitePool,
    ) -> impl Filter<Extract = (impl warp::Reply,), Error = warp::Rejection> + Clone {
        warp::path!("menu" / String)
            .and(warp::get())
            .and(with_db(pool.clone()))
            .and_then(|name: String, pool: SqlitePool| async move {
                let mut conn = pool.acquire().await.unwrap();
                let items = sqlx::query_as::<_, MenuItem>(
                    "SELECT id, name, price FROM MENU_ITEM WHERE menu = ?1",
                )
                .bind(name.clone())
                .fetch_all(&mut *conn)
                .await
                .unwrap();
                Ok::<warp::reply::Json, warp::Rejection>(warp::reply::json(&items))
            })
    }

    fn update_menu(
        pool: &SqlitePool,
    ) -> impl Filter<Extract = (impl warp::Reply,), Error = warp::Rejection> + Clone {
        warp::path!("menu" / String)
            .and(warp::put())
            .and(json_body())
            .and(with_db(pool.clone()))
            .and_then(|name: String, menu: Menu, pool: SqlitePool| async move {
                let mut conn = pool.acquire().await.unwrap();
                for it in menu.items {
                    sqlx::query("INSERT INTO MENU_ITEM VALUES (?1, ?2, ?3, ?4)")
                        .bind(name.clone())
                        .bind(it.id)
                        .bind(it.name)
                        .bind(it.price as i64)
                        .execute(&mut *conn)
                        .await
                        .unwrap();
                }
                Ok::<StatusCode, warp::Rejection>(StatusCode::OK)
            })
    }
}
