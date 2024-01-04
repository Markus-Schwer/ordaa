use db::Db;
use search::{init_search_index, SearchContextReader};
use warp::Filter;

mod menu;
mod search;
mod frontend;
mod db;

pub fn routes(
    db: Db,
    ctx: SearchContextReader,
) -> impl Filter<Extract = (impl warp::Reply,), Error = warp::Rejection> + Clone {
    menu::filters::all(db, ctx).or(frontend::filters::all())
}

#[tokio::main]
async fn main() {
    let (search_writer, search_reader) = init_search_index();
    let index_writer_handle = search_writer.start_index_writer(search_reader.clone());

    let db = Db::new().await;
    db.init_schema().await;
    let server_handle = warp::serve(routes(db, search_reader)).run(([127, 0, 0, 1], 8080));
    let (_, _) = tokio::join!(server_handle, index_writer_handle);
}

pub mod filters {
    use serde::de::DeserializeOwned;
    use warp::Filter;

    use crate::{search::SearchContextReader, db::Db};

    pub fn with_db(
        db: Db,
    ) -> impl Filter<Extract = (Db,), Error = std::convert::Infallible> + Clone {
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
