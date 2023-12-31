use sqlx::{sqlite::SqliteConnectOptions, SqlitePool};

mod menu;

#[tokio::main]
async fn main() {
    let options = SqliteConnectOptions::new()
        .create_if_missing(true)
        .filename("db.sqlite");
    let pool = SqlitePool::connect_with(options).await.unwrap();

    let mut conn = pool.acquire().await.unwrap();
    sqlx::query(
        "CREATE TABLE IF NOT EXISTS MENU_ITEM (menu TEXT, id TEXT, name TEXT, price INTEGER);",
    )
    .execute(&mut *conn)
    .await
    .unwrap();

    warp::serve(menu::filters::all(&pool))
        .run(([127, 0, 0, 1], 8080))
        .await;
}

pub mod filters {
    use serde::de::DeserializeOwned;
    use sqlx::SqlitePool;
    use warp::Filter;

    pub fn with_db(
        db: SqlitePool,
    ) -> impl Filter<Extract = (SqlitePool,), Error = std::convert::Infallible> + Clone {
        warp::any().map(move || db.clone())
    }

    pub fn json_body<T: Send + DeserializeOwned>(
    ) -> impl Filter<Extract = (T,), Error = warp::Rejection> + Clone {
        // When accepting a body, we want a JSON body
        // (and to reject huge payloads)...
        warp::body::content_length_limit(1024 * 16).and(warp::body::json())
    }
}
