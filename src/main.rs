use search::{init_search_index, SearchContextReader};
use warp::Filter;
use sqlx::{sqlite::SqliteConnectOptions, SqlitePool};

mod menu;
mod search;
mod frontend;

pub fn routes(
    pool: &SqlitePool,
    ctx: SearchContextReader,
) -> impl Filter<Extract = (impl warp::Reply,), Error = warp::Rejection> + Clone {
    menu::filters::all(pool, ctx).or(frontend::filters::all(pool))
}

#[tokio::main]
async fn main() {
    let options = SqliteConnectOptions::new()
        .create_if_missing(true)
        .filename("db.sqlite");
    let pool = SqlitePool::connect_with(options).await.unwrap();

    let (search_writer, search_reader) = init_search_index();
    let index_writer_handle = search_writer.start_index_writer(search_reader.clone());

    let mut conn = pool.acquire().await.unwrap();
    sqlx::query(
        "CREATE TABLE IF NOT EXISTS MENU_ITEM (menu TEXT, id TEXT, name TEXT, price INTEGER);",
    )
    .execute(&mut *conn)
    .await
    .unwrap();
    let server_handle = warp::serve(routes(&pool, search_reader)).run(([127, 0, 0, 1], 8080));
    let (_, _) = tokio::join!(server_handle, index_writer_handle);
}

pub mod filters {
    use serde::de::DeserializeOwned;
    use sqlx::SqlitePool;
    use warp::Filter;

    use crate::search::SearchContextReader;

    pub fn with_db(
        db: SqlitePool,
    ) -> impl Filter<Extract = (SqlitePool,), Error = std::convert::Infallible> + Clone {
        warp::any().map(move || db.clone())
    }

    pub fn with_searcher_ctx(
        ctx: SearchContextReader,
    ) -> impl Filter<Extract = (SearchContextReader,), Error = std::convert::Infallible> + Clone
    {
        warp::any().map(move || ctx.clone())
    }

    pub fn json_body<T: Send + DeserializeOwned>(
    ) -> impl Filter<Extract = (T,), Error = warp::Rejection> + Clone {
        // When accepting a body, we want a JSON body
        // (and to reject huge payloads)...
        warp::body::content_length_limit(1024 * 16).and(warp::body::json())
    }
}
