use std::path::Path;

use menu::Menu;
use sqlx::{sqlite::SqliteConnectOptions, SqlitePool};
use warp::Filter;
use tantivy::{query::QueryParser, schema::*, Index, IndexReader};
use tokio::sync::mpsc::{unbounded_channel, UnboundedSender};

mod menu;
mod frontend;

pub fn routes(
    pool: &SqlitePool,
    ctx: SearchContext,
) -> impl Filter<Extract = (impl warp::Reply,), Error = warp::Rejection> + Clone {
    menu::filters::all(pool, ctx).or(frontend::filters::all(pool))
}

#[derive(Clone)]
pub struct SearchContext {
    pub menu_field: Field,
    pub id_field: Field,
    pub name_field: Field,
    pub index_write_sender: UnboundedSender<Menu>,
    pub parser: QueryParser,
    pub reader: IndexReader,
}

#[tokio::main]
async fn main() {
    let options = SqliteConnectOptions::new()
        .create_if_missing(true)
        .filename("db.sqlite");
    let pool = SqlitePool::connect_with(options).await.unwrap();

    let (index_write_sender, mut index_write_receiver) = unbounded_channel::<Menu>();

    let mut schema_builder = Schema::builder();
    let id_field = schema_builder.add_text_field("id", TEXT | STORED);
    let menu_field = schema_builder.add_text_field("menu", TEXT | STORED);
    let name_field = schema_builder.add_text_field("name", TEXT | STORED);
    let index = match Index::open_in_dir(Path::new("./index")) {
        Ok(index) => index,
        Err(_) => Index::create_in_dir(Path::new("./index"), schema_builder.build()).unwrap(),
    };
    let ctx = SearchContext {
        index_write_sender,
        id_field,
        menu_field,
        name_field,
        reader: index.reader().unwrap(),
        parser: QueryParser::for_index(&index, vec![id_field, name_field]),
    };

    let index_writer_handle = tokio::spawn({
        let index = index.clone();
        async move {
            let mut writer = index.writer(500_000_000).unwrap();
            while let Some(menu_val) = index_write_receiver.recv().await {
                for it in menu_val.items {
                    writer
                        .add_document(tantivy::doc!(
                            ctx.menu_field => menu_val.name.clone(),
                            ctx.id_field => it.id,
                            ctx.name_field => it.name
                        ))
                        .unwrap();
                    writer.commit().unwrap();
                }
            }
        }
    });

    let mut conn = pool.acquire().await.unwrap();
    sqlx::query(
        "CREATE TABLE IF NOT EXISTS MENU_ITEM (menu TEXT, id TEXT, name TEXT, price INTEGER);",
    )
    .execute(&mut *conn)
    .await
    .unwrap();
    let server_handle = warp::serve(routes(&pool, ctx)).run(([127, 0, 0, 1], 8080));
    let (_, _) = tokio::join!(server_handle, index_writer_handle);
}

pub mod filters {
    use serde::de::DeserializeOwned;
    use sqlx::SqlitePool;
    use warp::Filter;

    use crate::SearchContext;

    pub fn with_db(
        db: SqlitePool,
    ) -> impl Filter<Extract = (SqlitePool,), Error = std::convert::Infallible> + Clone {
        warp::any().map(move || db.clone())
    }

    pub fn with_searcher_ctx(
        ctx: SearchContext,
    ) -> impl Filter<Extract = (SearchContext,), Error = std::convert::Infallible> + Clone {
        warp::any().map(move || ctx.clone())
    }

    pub fn json_body<T: Send + DeserializeOwned>(
    ) -> impl Filter<Extract = (T,), Error = warp::Rejection> + Clone {
        // When accepting a body, we want a JSON body
        // (and to reject huge payloads)...
        warp::body::content_length_limit(1024 * 16).and(warp::body::json())
    }
}
